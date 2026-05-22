package license

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
)

func TestGenerateLicenseSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll failed: %v", err)
		}
		if !strings.Contains(string(body), `"provider":"https://provider.example.test"`) {
			t.Fatalf("expected provider in payload, got %s", string(body))
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"lic-generated"}`))
	}))
	defer server.Close()

	service := NewService(server.URL, "user", "pass", "", "", "", "https://provider.example.test")
	start := time.Unix(100, 0).UTC()
	end := start.Add(24 * time.Hour)
	rightPrint := 10
	rightCopy := 5
	license := &domain.License{
		PublicationID: "pub-1",
		UserID:        "reader-1",
		Passphrase:    "secret",
		Hint:          "pet",
		StartDate:     &start,
		EndDate:       &end,
		RightPrint:    &rightPrint,
		RightCopy:     &rightCopy,
	}

	if err := service.GenerateLicense(context.Background(), license); err != nil {
		t.Fatalf("GenerateLicense failed: %v", err)
	}
	if license.ID != "lic-generated" {
		t.Fatalf("expected generated license id, got %q", license.ID)
	}
	if !strings.Contains(license.LCPL, `"id":"lic-generated"`) {
		t.Fatalf("expected lcpl payload to be stored, got %s", license.LCPL)
	}
}

func TestGenerateLicenseNotFoundAndValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	service := NewService(server.URL, "", "", "", "", "", "https://provider.example.test")
	err := service.GenerateLicense(context.Background(), &domain.License{PublicationID: "pub-1", UserID: "reader-1"})
	if err != ErrContentNotFound {
		t.Fatalf("expected ErrContentNotFound, got %v", err)
	}

	if err := service.GenerateLicense(context.Background(), &domain.License{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetLicenseAndHash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".lcpl") || strings.HasSuffix(r.URL.Path, "/lic-1") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"lic-1"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	service := NewService(server.URL, "", "", "", "", "", "")
	body, err := service.GetLicense(context.Background(), "lic-1")
	if err != nil {
		t.Fatalf("GetLicense failed: %v", err)
	}
	if !strings.Contains(string(body), `"id":"lic-1"`) {
		t.Fatalf("unexpected license body %s", string(body))
	}
	if lcpPassphraseHash("secret") == "" {
		t.Fatal("expected passphrase hash")
	}
}
