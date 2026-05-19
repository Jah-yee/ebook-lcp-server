package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
	usecaseLicense "github.com/amirhdev/ebook-lcp-server/internal/usecase/lcp/license"
)

type AdminLicensesHandler struct {
	licenses usecaseLicense.LicenseUsecase
}

type licensePatchRequest struct {
	EndDate *string `json:"endDate"`
}

type adminLicenseResponse struct {
	ID             string     `json:"id"`
	PublicationID  string     `json:"publicationID"`
	UserID         string     `json:"userID"`
	Passphrase     string     `json:"passphrase"`
	Hint           string     `json:"hint"`
	PublicationURL string     `json:"publicationURL"`
	RightPrint     *int       `json:"rightPrint,omitempty"`
	RightCopy      *int       `json:"rightCopy,omitempty"`
	StartDate      *time.Time `json:"startDate,omitempty"`
	EndDate        *time.Time `json:"endDate,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"createdAt"`
}

func NewAdminLicensesHandler(licenses usecaseLicense.LicenseUsecase) *AdminLicensesHandler {
	return &AdminLicensesHandler{licenses: licenses}
}

func (h *AdminLicensesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/admin/licenses"), "/")
	if path == "" {
		h.list(w, r)
		return
	}

	parts := strings.Split(path, "/")
	licenseID := parts[0]
	if len(parts) == 2 && parts[1] == "revoke" {
		h.revoke(w, r, licenseID)
		return
	}
	if len(parts) == 1 {
		h.update(w, r, licenseID)
		return
	}
	http.NotFound(w, r)
}

func (h *AdminLicensesHandler) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	licenses, err := h.licenses.GetByPublication(r.Context(), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	items := make([]adminLicenseResponse, 0, len(licenses))
	for _, lic := range licenses {
		items = append(items, toAdminLicenseResponse(lic))
	}
	writeJSON(w, http.StatusOK, map[string]any{"licenses": items})
}

func (h *AdminLicensesHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req licensePatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	var endDate *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsed, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "endDate must be RFC3339"})
			return
		}
		endDate = &parsed
	}
	lic, err := h.licenses.UpdateEndDate(r.Context(), id, endDate)
	if err != nil {
		if err.Error() == "license not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, toAdminLicenseResponse(lic))
}

func (h *AdminLicensesHandler) revoke(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := h.licenses.Revoke(r.Context(), id); err != nil {
		if err.Error() == "license not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func toAdminLicenseResponse(lic *lcp.License) adminLicenseResponse {
	status := lic.Status
	if status == "" {
		status = "ready"
	}
	return adminLicenseResponse{
		ID:             lic.ID,
		PublicationID:  lic.PublicationID,
		UserID:         lic.UserID,
		Passphrase:     lic.Passphrase,
		Hint:           lic.Hint,
		PublicationURL: lic.PublicationURL,
		RightPrint:     lic.RightPrint,
		RightCopy:      lic.RightCopy,
		StartDate:      lic.StartDate,
		EndDate:        lic.EndDate,
		Status:         status,
		CreatedAt:      lic.CreatedAt,
	}
}
