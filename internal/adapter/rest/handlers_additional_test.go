package rest

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/amirhdev/ebook-lcp-server/internal/auth"
	domain "github.com/amirhdev/ebook-lcp-server/internal/domain"
	"github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
)

type extraLicenseUsecase struct {
	licenses []*lcp.License
	err      error
}

func (u *extraLicenseUsecase) Create(context.Context, *lcp.LicenseInput) (*lcp.License, error) {
	return nil, nil
}

func (u *extraLicenseUsecase) GetByID(_ context.Context, id string) (*lcp.License, error) {
	if u.err != nil {
		return nil, u.err
	}
	for _, lic := range u.licenses {
		if lic.ID == id {
			return lic, nil
		}
	}
	return nil, nil
}

func (u *extraLicenseUsecase) GetByPublication(context.Context, *string) ([]*lcp.License, error) {
	return u.licenses, nil
}

func (u *extraLicenseUsecase) UpdateEndDate(context.Context, string, *time.Time) (*lcp.License, error) {
	return nil, nil
}

func (u *extraLicenseUsecase) Revoke(context.Context, string) error {
	return nil
}

type extraPublicationRepo struct {
	items map[string]*lcp.Publication
	err   error
}

func (r *extraPublicationRepo) Save(_ context.Context, pub *lcp.Publication) error {
	if r.err != nil {
		return r.err
	}
	if r.items == nil {
		r.items = map[string]*lcp.Publication{}
	}
	copy := *pub
	r.items[pub.ID] = &copy
	return nil
}

func (r *extraPublicationRepo) FindAll(context.Context) ([]*lcp.Publication, error) {
	if r.err != nil {
		return nil, r.err
	}
	items := make([]*lcp.Publication, 0, len(r.items))
	for _, item := range r.items {
		copy := *item
		items = append(items, &copy)
	}
	return items, nil
}

func (r *extraPublicationRepo) FindByID(_ context.Context, id string) (*lcp.Publication, error) {
	if r.err != nil {
		return nil, r.err
	}
	item, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	copy := *item
	return &copy, nil
}

type extraAuditRepo struct {
	entries []*domain.AuditEntry
	err     error
}

func (r *extraAuditRepo) Save(context.Context, *domain.AuditEntry) error {
	return nil
}

func (r *extraAuditRepo) FindRecent(context.Context, int) ([]*domain.AuditEntry, error) {
	return r.entries, r.err
}

func (r *extraAuditRepo) FindRecentByTenant(_ context.Context, tenantID string, _ int) ([]*domain.AuditEntry, error) {
	if r.err != nil {
		return nil, r.err
	}
	items := make([]*domain.AuditEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		if entry.TenantID == tenantID {
			items = append(items, entry)
		}
	}
	return items, nil
}

type noopLCPLProvider struct{}

func (noopLCPLProvider) GetLicense(context.Context, string) ([]byte, error) {
	return nil, nil
}

func TestStaticFileHandlerServesAndRejectsMethods(t *testing.T) {
	path := filepath.Join(t.TempDir(), "swagger.yaml")
	if err := os.WriteFile(path, []byte("openapi: 3.0.0"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	handler := staticFileHandler(path, "text/yaml; charset=utf-8")
	getReq := httptest.NewRequest(http.MethodGet, "/swagger.yaml", nil)
	getRec := httptest.NewRecorder()
	handler(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected swagger get to succeed, got %d", getRec.Code)
	}
	if !strings.Contains(getRec.Header().Get("Content-Type"), "text/yaml") {
		t.Fatalf("unexpected content type %q", getRec.Header().Get("Content-Type"))
	}

	postReq := httptest.NewRequest(http.MethodPost, "/swagger.yaml", nil)
	postRec := httptest.NewRecorder()
	handler(postRec, postReq)
	if postRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed, got %d", postRec.Code)
	}
}

func TestLicenseDownloadHandler(t *testing.T) {
	usecase := &extraLicenseUsecase{licenses: []*lcp.License{{ID: "lic-1", LCPL: `{"id":"lic-1"}`}}}
	handler := NewLicenseDownloadHandler(usecase, noopLCPLProvider{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/licenses/lic-1.lcpl", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "lic-1") {
		t.Fatalf("unexpected download response %d: %s", rec.Code, rec.Body.String())
	}

	notFound := httptest.NewRecorder()
	handler.ServeHTTP(notFound, httptest.NewRequest(http.MethodGet, "/licenses/missing.lcpl", nil))
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("expected missing license to return 404, got %d", notFound.Code)
	}
}

func TestExtractLicenseIDForLCPL(t *testing.T) {
	cases := map[string]string{
		"/licenses/lic-1.lcpl":        "lic-1",
		"/api/v1/licenses/lic-2/lcpl": "lic-2",
		"licenses/lic-3.lcpl":         "lic-3",
		"api/v1/licenses/lic-4/lcpl":  "lic-4",
	}
	for input, want := range cases {
		got, ok := extractLicenseIDForLCPL(input)
		if !ok || got != want {
			t.Fatalf("expected %q from %q, got %q ok=%v", want, input, got, ok)
		}
	}
	if _, ok := extractLicenseIDForLCPL("/licenses"); ok {
		t.Fatal("expected invalid path to fail")
	}
}

func TestLicenseStatusDocument(t *testing.T) {
	usecase := &extraLicenseUsecase{licenses: []*lcp.License{{ID: "lic-1"}}}
	handler := LicenseStatusDocument(usecase)
	req := httptest.NewRequest(http.MethodGet, "/licenses/lic-1/status", nil)
	req.Host = "example.test"
	req.Header.Set("X-Forwarded-Proto", "http")
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ready"`) {
		t.Fatalf("expected default ready status, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "http://example.test/licenses/lic-1.lcpl") {
		t.Fatalf("expected license link in response, got %s", rec.Body.String())
	}
}

func TestLicenseUserData(t *testing.T) {
	usecase := &extraLicenseUsecase{licenses: []*lcp.License{{ID: "lic-1", UserID: "reader-1", Passphrase: "secret", Hint: "pet"}}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/licenses/lic-1/user", nil)
	LicenseUserData(usecase)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	wantSum := sha256.Sum256([]byte("secret"))
	wantHash := hex.EncodeToString(wantSum[:])
	if !strings.Contains(rec.Body.String(), wantHash) {
		t.Fatalf("expected passphrase hash in response, got %s", rec.Body.String())
	}
}

func TestAuditHandlerFiltersByTenant(t *testing.T) {
	handler := NewAuditHandler(&extraAuditRepo{entries: []*domain.AuditEntry{
		{ID: "1", TenantID: "tenant-a"},
		{ID: "2", TenantID: "tenant-b"},
	}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit?limit=10", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), &auth.Claims{TenantID: "tenant-a"}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"1"`) || strings.Contains(rec.Body.String(), `"2"`) {
		t.Fatalf("unexpected audit response %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAuditHandlerReturnsInternalError(t *testing.T) {
	handler := NewAuditHandler(&extraAuditRepo{err: errors.New("boom")})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected internal error, got %d", rec.Code)
	}
}

func TestAdminUserStoreAndHandler(t *testing.T) {
	store := NewAdminUserStore(t.TempDir())
	users := store.List()
	if len(users) == 0 {
		t.Fatal("expected default users")
	}
	updated, err := store.SetVerified(users[0].ID, false)
	if err != nil || updated == nil || updated.Verified {
		t.Fatalf("expected user verification to update, got user=%+v err=%v", updated, err)
	}

	handler := NewAdminUsersHandler(store)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestAdminUserStoreLoadFallbacks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	store := NewAdminUserStore(dir)
	if len(store.List()) == 0 {
		t.Fatal("expected default users after invalid persisted data")
	}
}

func TestAdminUsersHandlerToggleRoutes(t *testing.T) {
	store := NewAdminUserStore(t.TempDir())
	users := store.List()
	handler := NewAdminUsersHandler(store)

	verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/"+users[1].ID+"/verify", nil)
	verifyRec := httptest.NewRecorder()
	handler.ServeHTTP(verifyRec, verifyReq)
	if verifyRec.Code != http.StatusOK {
		t.Fatalf("expected verify response 200, got %d", verifyRec.Code)
	}

	unverifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/"+users[0].ID+"/unverify", nil)
	unverifyRec := httptest.NewRecorder()
	handler.ServeHTTP(unverifyRec, unverifyReq)
	if unverifyRec.Code != http.StatusOK {
		t.Fatalf("expected unverify response 200, got %d", unverifyRec.Code)
	}

	missingReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/missing/verify", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("expected missing user response 404, got %d", missingRec.Code)
	}
}

func TestAdminTenantsHandlerAdditionalBranches(t *testing.T) {
	store := NewTenantStore(t.TempDir(), "default")
	handler := NewAdminTenantsHandler(store)

	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/admin/tenants/tenant-a", strings.NewReader(`{"name":"Tenant A"}`))
	putRec := httptest.NewRecorder()
	handler.ServeHTTP(putRec, putReq)
	if putRec.Code != http.StatusOK {
		t.Fatalf("expected put response 200, got %d", putRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/tenants/tenant-a", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected get response 200, got %d", getRec.Code)
	}

	methodReq := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/tenants", nil)
	methodRec := httptest.NewRecorder()
	handler.ServeHTTP(methodRec, methodReq)
	if methodRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed, got %d", methodRec.Code)
	}

	badJSONReq := httptest.NewRequest(http.MethodPut, "/api/v1/admin/tenants/tenant-b", strings.NewReader("{"))
	badJSONRec := httptest.NewRecorder()
	handler.ServeHTTP(badJSONRec, badJSONReq)
	if badJSONRec.Code != http.StatusBadRequest {
		t.Fatalf("expected bad json response 400, got %d", badJSONRec.Code)
	}
}

func TestPublicationHandlerCreateListPatchAndStatus(t *testing.T) {
	repo := &extraPublicationRepo{items: map[string]*lcp.Publication{
		"pub-2": {ID: "pub-2", Title: "Other", TenantID: "tenant-b"},
	}}
	handler := NewPublicationHandler(repo, fakePublicationUsecase{})

	body := `{"title":"Book","encrypted_uri":"s3://bucket/book.epub","right_print":5}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/publications", strings.NewReader(body))
	req = req.WithContext(auth.WithClaims(req.Context(), &auth.Claims{Role: "publisher", Roles: []string{"publisher"}, TenantID: "tenant-a"}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected create 201, got %d: %s", rec.Code, rec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/publications", nil)
	listReq = listReq.WithContext(auth.WithClaims(listReq.Context(), &auth.Claims{TenantID: "tenant-a"}))
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK || strings.Contains(listRec.Body.String(), `"pub-2"`) {
		t.Fatalf("unexpected list response %d: %s", listRec.Code, listRec.Body.String())
	}

	var pubID string
	for id := range repo.items {
		if id != "pub-2" {
			pubID = id
		}
	}
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/publications/"+pubID, strings.NewReader(`{"title":"Renamed","right_copy":3}`))
	patchReq = patchReq.WithContext(auth.WithClaims(patchReq.Context(), &auth.Claims{Role: "admin", Roles: []string{"admin"}, TenantID: "tenant-a"}))
	patchRec := httptest.NewRecorder()
	handler.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusOK || repo.items[pubID].Title != "Renamed" {
		t.Fatalf("unexpected patch response %d: %s", patchRec.Code, patchRec.Body.String())
	}

	statusReq := httptest.NewRequest(http.MethodPost, "/api/v1/publications/"+pubID+"/deactivate", nil)
	statusReq = statusReq.WithContext(auth.WithClaims(statusReq.Context(), &auth.Claims{Role: "publisher", Roles: []string{"publisher"}, TenantID: "tenant-a"}))
	statusRec := httptest.NewRecorder()
	handler.ServeHTTP(statusRec, statusReq)
	if statusRec.Code != http.StatusOK || repo.items[pubID].Status != "inactive" {
		t.Fatalf("unexpected deactivate response %d: %s", statusRec.Code, statusRec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/publications/"+pubID, nil)
	getReq = getReq.WithContext(auth.WithClaims(getReq.Context(), &auth.Claims{TenantID: "tenant-a"}))
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK || !strings.Contains(getRec.Body.String(), `"Renamed"`) {
		t.Fatalf("unexpected get response %d: %s", getRec.Code, getRec.Body.String())
	}

	hiddenReq := httptest.NewRequest(http.MethodGet, "/api/v1/publications/pub-2", nil)
	hiddenReq = hiddenReq.WithContext(auth.WithClaims(hiddenReq.Context(), &auth.Claims{TenantID: "tenant-a"}))
	hiddenRec := httptest.NewRecorder()
	handler.ServeHTTP(hiddenRec, hiddenReq)
	if hiddenRec.Code != http.StatusNotFound {
		t.Fatalf("expected hidden tenant publication to return 404, got %d", hiddenRec.Code)
	}
}

func TestPublicationHandlerRejectsBadInputAndUnauthorized(t *testing.T) {
	repo := &extraPublicationRepo{items: map[string]*lcp.Publication{}}
	handler := NewPublicationHandler(repo, fakePublicationUsecase{})

	badCreate := httptest.NewRequest(http.MethodPost, "/api/v1/publications", strings.NewReader(`{"title":"Book","file":"!!!"}`))
	badCreate = badCreate.WithContext(auth.WithClaims(badCreate.Context(), &auth.Claims{Role: "publisher", Roles: []string{"publisher"}}))
	badRec := httptest.NewRecorder()
	handler.ServeHTTP(badRec, badCreate)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for invalid base64, got %d", badRec.Code)
	}

	unauthReq := httptest.NewRequest(http.MethodPost, "/api/v1/publications", strings.NewReader(`{"title":"Book","encrypted_uri":"s3://bucket/book.epub"}`))
	unauthRec := httptest.NewRecorder()
	handler.ServeHTTP(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden response, got %d", unauthRec.Code)
	}

	negative := -1
	if err := validatePublicationRights(&negative, nil); err == nil {
		t.Fatal("expected validation error")
	}
	if defaultPublicationStatus("") != "active" || defaultPublicationStatus(" INACTIVE ") != "inactive" {
		t.Fatal("expected normalized publication status")
	}
	if !claimsHasRole(&auth.Claims{Roles: []string{"publisher"}}, "publisher") {
		t.Fatal("expected claimsHasRole to match roles slice")
	}

	missingReq := httptest.NewRequest(http.MethodGet, "/api/v1/publications/missing", nil)
	missingReq = missingReq.WithContext(auth.WithClaims(missingReq.Context(), &auth.Claims{TenantID: "tenant-a"}))
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("expected missing publication to return 404, got %d", missingRec.Code)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/publications/missing", strings.NewReader("{"))
	patchReq = patchReq.WithContext(auth.WithClaims(patchReq.Context(), &auth.Claims{Role: "admin", Roles: []string{"admin"}, TenantID: "tenant-a"}))
	patchRec := httptest.NewRecorder()
	handler.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusNotFound {
		t.Fatalf("expected patch missing publication to return 404, got %d", patchRec.Code)
	}
}

func TestPublicationHandlerUploadPath(t *testing.T) {
	repo := &extraPublicationRepo{items: map[string]*lcp.Publication{}}
	handler := NewPublicationHandler(repo, fakePublicationUsecase{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/publications", strings.NewReader(`{"title":"Book","file":"`+base64.StdEncoding.EncodeToString([]byte("hello"))+`"}`))
	req = req.WithContext(auth.WithClaims(req.Context(), &auth.Claims{Role: "publisher", Roles: []string{"publisher"}, TenantID: "tenant-a"}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected upload create 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLCPPassphraseHash(t *testing.T) {
	sum := sha256.Sum256([]byte("secret"))
	if got := lcpPassphraseHash("secret"); got != hex.EncodeToString(sum[:]) {
		t.Fatalf("unexpected passphrase hash %q", got)
	}
}

func TestLicenseDownloadRejectsMethod(t *testing.T) {
	handler := NewLicenseDownloadHandler(&extraLicenseUsecase{}, noopLCPLProvider{})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/licenses/lic-1.lcpl", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed, got %d", rec.Code)
	}
}

func TestLicenseStatusDocumentRejectsMissingLicense(t *testing.T) {
	handler := LicenseStatusDocument(&extraLicenseUsecase{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/licenses/missing/status", nil)
	handler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", rec.Code)
	}
}

func TestLicenseUserDataRejectsMissingLicense(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/licenses/missing/user", nil)
	LicenseUserData(&extraLicenseUsecase{})(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", rec.Code)
	}
}
