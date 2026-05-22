package audit

import (
	"context"
	"testing"

	repo "github.com/amirhdev/ebook-lcp-server/internal/adapter/repository/audit"
	"github.com/amirhdev/ebook-lcp-server/internal/auth"
	domain "github.com/amirhdev/ebook-lcp-server/internal/domain"
)

type stubRepository struct {
	saved *domain.AuditEntry
}

func (r *stubRepository) Save(_ context.Context, entry *domain.AuditEntry) error {
	r.saved = entry
	return nil
}

func (r *stubRepository) FindRecent(context.Context, int) ([]*domain.AuditEntry, error) {
	return nil, nil
}

func (r *stubRepository) FindRecentByTenant(context.Context, string, int) ([]*domain.AuditEntry, error) {
	return nil, nil
}

var _ repo.Repository = (*stubRepository)(nil)

func TestRecordUsesClaimsSubjectAndTenant(t *testing.T) {
	store := &stubRepository{}
	service := NewService(store)
	ctx := auth.WithClaims(context.Background(), &auth.Claims{Subject: "alice", TenantID: "tenant-a"})

	if err := service.Record(ctx, "license.created", "license", "lic-1"); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if store.saved == nil {
		t.Fatal("expected audit entry to be saved")
	}
	if store.saved.Actor != "alice" || store.saved.TenantID != "tenant-a" {
		t.Fatalf("unexpected saved audit entry: %+v", store.saved)
	}
}

func TestRecordDefaultsToSystemActor(t *testing.T) {
	store := &stubRepository{}
	service := NewService(store)

	if err := service.Record(context.Background(), "publication.uploaded", "publication", "pub-1"); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if store.saved.Actor != "system" {
		t.Fatalf("expected system actor, got %q", store.saved.Actor)
	}
}
