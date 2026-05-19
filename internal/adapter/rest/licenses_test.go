package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
)

type fakeAdminLicenseUsecase struct {
	licenses []*lcp.License
	updated  *time.Time
	revoked  string
}

func (f *fakeAdminLicenseUsecase) Create(context.Context, *lcp.LicenseInput) (*lcp.License, error) {
	return nil, nil
}

func (f *fakeAdminLicenseUsecase) GetByID(_ context.Context, id string) (*lcp.License, error) {
	for _, lic := range f.licenses {
		if lic.ID == id {
			return lic, nil
		}
	}
	return nil, nil
}

func (f *fakeAdminLicenseUsecase) GetByPublication(context.Context, *string) ([]*lcp.License, error) {
	return f.licenses, nil
}

func (f *fakeAdminLicenseUsecase) UpdateEndDate(_ context.Context, id string, endDate *time.Time) (*lcp.License, error) {
	f.updated = endDate
	lic, _ := f.GetByID(context.Background(), id)
	lic.EndDate = endDate
	return lic, nil
}

func (f *fakeAdminLicenseUsecase) Revoke(_ context.Context, id string) error {
	f.revoked = id
	return nil
}

func TestAdminLicensesList(t *testing.T) {
	handler := NewAdminLicensesHandler(&fakeAdminLicenseUsecase{
		licenses: []*lcp.License{{ID: "lic1", Status: "ready"}},
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/licenses", nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"lic1"`) {
		t.Fatalf("unexpected response %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAdminLicensesUpdateAndRevoke(t *testing.T) {
	usecase := &fakeAdminLicenseUsecase{licenses: []*lcp.License{{ID: "lic1", Status: "ready"}}}
	handler := NewAdminLicensesHandler(usecase)

	patchRec := httptest.NewRecorder()
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/licenses/lic1", strings.NewReader(`{"endDate":"2026-05-31T00:00:00Z"}`))
	handler.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusOK || usecase.updated == nil {
		t.Fatalf("unexpected patch response %d: %s", patchRec.Code, patchRec.Body.String())
	}

	revokeRec := httptest.NewRecorder()
	revokeReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/licenses/lic1/revoke", nil)
	handler.ServeHTTP(revokeRec, revokeReq)
	if revokeRec.Code != http.StatusOK || usecase.revoked != "lic1" {
		t.Fatalf("unexpected revoke response %d: %s", revokeRec.Code, revokeRec.Body.String())
	}
}
