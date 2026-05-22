package audit

import (
	"context"
	"path/filepath"
	"testing"

	domain "github.com/amirhdev/ebook-lcp-server/internal/domain"
)

func TestRepositorySaveAndFindRecent(t *testing.T) {
	repo := NewRepository()
	entry1 := &domain.AuditEntry{ID: "1", TenantID: "tenant-a", Action: "one"}
	entry2 := &domain.AuditEntry{ID: "2", TenantID: "tenant-b", Action: "two"}
	if err := repo.Save(context.Background(), entry1); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if err := repo.Save(context.Background(), entry2); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	items, err := repo.FindRecent(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindRecent failed: %v", err)
	}
	if len(items) != 1 || items[0].ID != "2" {
		t.Fatalf("unexpected recent entries: %+v", items)
	}

	tenantItems, err := repo.FindRecentByTenant(context.Background(), "tenant-a", 10)
	if err != nil {
		t.Fatalf("FindRecentByTenant failed: %v", err)
	}
	if len(tenantItems) != 1 || tenantItems[0].ID != "1" {
		t.Fatalf("unexpected tenant entries: %+v", tenantItems)
	}
}

func TestPersistentRepositoryLoadsSavedEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.json")
	repo, err := NewPersistentRepository(path)
	if err != nil {
		t.Fatalf("NewPersistentRepository failed: %v", err)
	}
	if err := repo.Save(context.Background(), &domain.AuditEntry{ID: "1", TenantID: "tenant-a"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	reloaded, err := NewPersistentRepository(path)
	if err != nil {
		t.Fatalf("NewPersistentRepository failed: %v", err)
	}
	items, err := reloaded.FindRecent(context.Background(), 10)
	if err != nil {
		t.Fatalf("FindRecent failed: %v", err)
	}
	if len(items) != 1 || items[0].ID != "1" {
		t.Fatalf("unexpected persisted entries: %+v", items)
	}
}

func TestSanitizeAuditLimit(t *testing.T) {
	if got := sanitizeAuditLimit(0); got != maxAuditEntriesLimit {
		t.Fatalf("expected max limit for zero input, got %d", got)
	}
	if got := sanitizeAuditLimit(maxAuditEntriesLimit + 1); got != maxAuditEntriesLimit {
		t.Fatalf("expected capped limit, got %d", got)
	}
	if got := sanitizeAuditLimit(5); got != 5 {
		t.Fatalf("expected explicit limit, got %d", got)
	}
}
