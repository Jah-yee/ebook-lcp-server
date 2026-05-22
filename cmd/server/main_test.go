package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/amirhdev/ebook-lcp-server/internal/adapter/rest"
	"github.com/amirhdev/ebook-lcp-server/internal/auth"
	"github.com/amirhdev/ebook-lcp-server/internal/config"
	domain "github.com/amirhdev/ebook-lcp-server/internal/domain"
	"github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
	publicationstorage "github.com/amirhdev/ebook-lcp-server/internal/storage"
	"github.com/amirhdev/ebook-lcp-server/internal/webhook"
)

type fakePublicationUsecase struct {
	pub *lcp.Publication
}

func (u fakePublicationUsecase) UploadAndEncrypt(context.Context, string, io.Reader) (*lcp.Publication, error) {
	return nil, nil
}

func (u fakePublicationUsecase) GetAll(context.Context) ([]*lcp.Publication, error) {
	return []*lcp.Publication{u.pub}, nil
}

func (u fakePublicationUsecase) GetByID(context.Context, string) (*lcp.Publication, error) {
	return u.pub, nil
}

type fakeSignedStorage struct{}

func (fakeSignedStorage) StoreEncrypted(context.Context, string, string) (string, error) {
	return "", nil
}

func (fakeSignedStorage) OpenEncrypted(context.Context, string) (io.ReadCloser, error) {
	return nil, nil
}

func (fakeSignedStorage) SignedURL(context.Context, string, time.Duration) (string, bool, error) {
	return "http://localhost:9000/books/publications/book.epub?signature=ok", true, nil
}

func (fakeSignedStorage) Ready(context.Context) error {
	return nil
}

func TestPublicationDownloadHandlerRedirectsToSignedURL(t *testing.T) {
	cfg := &config.Config{}
	cfg.LCP.Storage.S3.SignedURLTTLSecs = 900
	handler := publicationDownloadHandler(fakePublicationUsecase{
		pub: &lcp.Publication{
			ID:           "book",
			EncryptedURI: "s3://books/publications/book.epub",
		},
	}, fakeSignedStorage{}, cfg)

	req := httptest.NewRequest(http.MethodGet, "/publications/book/content", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}
	location := rec.Header().Get("Location")
	if !strings.Contains(location, "signature=ok") {
		t.Fatalf("unexpected redirect location: %s", location)
	}
}

func TestPublicationDownloadHandlerAcceptsEpubPathFromLCPL(t *testing.T) {
	cfg := &config.Config{}
	cfg.LCP.Storage.S3.SignedURLTTLSecs = 900
	handler := publicationDownloadHandler(fakePublicationUsecase{
		pub: &lcp.Publication{
			ID:           "book",
			EncryptedURI: "s3://books/publications/book.epub",
		},
	}, fakeSignedStorage{}, cfg)

	req := httptest.NewRequest(http.MethodGet, "/publications/book.epub", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}
}

type fakeStreamingStorage struct {
	body string
}

func (fakeStreamingStorage) StoreEncrypted(context.Context, string, string) (string, error) {
	return "", nil
}

func (s fakeStreamingStorage) OpenEncrypted(context.Context, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(s.body)), nil
}

func (fakeStreamingStorage) SignedURL(context.Context, string, time.Duration) (string, bool, error) {
	return "", false, nil
}

func (fakeStreamingStorage) Ready(context.Context) error {
	return nil
}

func TestPublicationDownloadHandlerStreamsWhenSignedURLDisabled(t *testing.T) {
	cfg := &config.Config{}
	handler := publicationDownloadHandler(fakePublicationUsecase{
		pub: &lcp.Publication{
			ID:           "book",
			EncryptedURI: "s3://books/publications/book.epub",
		},
	}, fakeStreamingStorage{body: "encrypted-book"}, cfg)

	req := httptest.NewRequest(http.MethodGet, "/publications/book/content", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "encrypted-book" {
		t.Fatalf("unexpected streamed body %q", rec.Body.String())
	}
}

func TestBuildBaseURLs(t *testing.T) {
	cfg := &config.Config{}
	if got := buildBaseURL(cfg); got != "http://localhost:8080" {
		t.Fatalf("unexpected default base URL %q", got)
	}

	cfg.Server.Port = ":9090"
	cfg.Server.PublicBaseURL = "https://example.test/base/"
	cfg.Server.StatusBaseURL = "https://status.example.test/"
	cfg.LCP.ProviderURI = "https://provider.example.test/"

	if got := buildBaseURL(cfg); got != "https://example.test/base" {
		t.Fatalf("unexpected public base URL %q", got)
	}
	if got := buildProviderURI(cfg); got != "https://provider.example.test" {
		t.Fatalf("unexpected provider URI %q", got)
	}
	if got := buildStatusBaseURL(cfg); got != "https://status.example.test" {
		t.Fatalf("unexpected status base URL %q", got)
	}
}

func TestTenantHelpersAndResolvers(t *testing.T) {
	store := rest.NewTenantStore(t.TempDir(), "default")
	if err := store.Upsert(&rest.TenantRecord{
		ID:           "tenant-a",
		Name:         "Tenant A",
		RateLimitRPM: 42,
		APIKeys: []rest.APIKey{{
			Key:     "api-key-1",
			Subject: "publisher-a",
			Role:    "publisher",
		}},
	}); err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	ctx := auth.WithClaims(context.Background(), &auth.Claims{TenantID: "tenant-a"})
	if got := tenantIDFromContext(ctx); got != "tenant-a" {
		t.Fatalf("unexpected tenant ID %q", got)
	}
	if got := tenantIDFromContext(context.Background()); got != "default" {
		t.Fatalf("unexpected default tenant ID %q", got)
	}

	resolveAPIKey := buildAPIKeyResolver(store)
	claims, ok := resolveAPIKey("api-key-1")
	if !ok || claims.Subject != "publisher-a" || claims.TenantID != "tenant-a" {
		t.Fatalf("unexpected API key claims %+v ok=%v", claims, ok)
	}
	if _, ok := resolveAPIKey("missing"); ok {
		t.Fatal("expected missing API key to fail")
	}

	resolveRateLimit := buildTenantRateLimitResolver(store)
	if got := resolveRateLimit("tenant-a"); got != 42 {
		t.Fatalf("unexpected tenant rate limit %d", got)
	}
	if got := resolveRateLimit("missing"); got != 0 {
		t.Fatalf("expected zero rate limit for missing tenant, got %d", got)
	}
}

func TestBuildReadyCheck(t *testing.T) {
	check := buildReadyCheck(nil, fakeSignedStorage{})
	if err := check(context.Background()); err != nil {
		t.Fatalf("expected ready check to pass, got %v", err)
	}
}

func TestBuildRepositoriesAndFailureRecorder(t *testing.T) {
	cfg := &config.Config{}
	cfg.DataDir = t.TempDir()

	pubRepo, err := buildPublicationRepository(cfg, nil)
	if err != nil {
		t.Fatalf("buildPublicationRepository failed: %v", err)
	}
	if err := pubRepo.Save(context.Background(), &lcp.Publication{ID: "pub-1"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	licRepo, err := buildLicenseRepository(cfg, nil)
	if err != nil {
		t.Fatalf("buildLicenseRepository failed: %v", err)
	}
	if err := licRepo.Save(context.Background(), &lcp.License{ID: "lic-1"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	auditRepo, err := buildAuditRepository(cfg, nil)
	if err != nil {
		t.Fatalf("buildAuditRepository failed: %v", err)
	}
	if err := auditRepo.Save(context.Background(), &domain.AuditEntry{ID: "audit-1"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	recorder := buildWebhookFailureRecorder(cfg)
	if err := recorder.Record(context.Background(), webhook.Failure{URL: "https://example.test", EventType: webhook.EventLicenseCreated}); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(cfg.DataDir, "webhook-failures.json"))
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	var failures []webhook.Failure
	if err := json.Unmarshal(raw, &failures); err != nil || len(failures) != 1 {
		t.Fatalf("unexpected failure file contents: err=%v failures=%+v", err, failures)
	}
}

func TestBuildPublicationStorageUsesFilesystemByDefault(t *testing.T) {
	cfg := &config.Config{}
	store, err := buildPublicationStorage(cfg, rest.NewTenantStore(t.TempDir(), "default"))
	if err != nil {
		t.Fatalf("buildPublicationStorage failed: %v", err)
	}
	if _, ok := store.(*publicationstorage.FilesystemPublicationStorage); !ok {
		t.Fatalf("expected filesystem storage, got %T", store)
	}
}

func TestBuildDatabaseReturnsNilWithoutDSN(t *testing.T) {
	db, err := buildDatabase(&config.Config{})
	if err != nil {
		t.Fatalf("buildDatabase failed: %v", err)
	}
	if db != nil {
		t.Fatalf("expected nil db, got %v", db)
	}
}
