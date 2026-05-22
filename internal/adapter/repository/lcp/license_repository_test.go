package lcp

import (
	"context"
	"path/filepath"
	"testing"

	domain "github.com/amirhdev/ebook-lcp-server/internal/domain/lcp"
)

func TestLicenseRepositorySaveFindDelete(t *testing.T) {
	repo := NewLicenseRepository()

	lic := &domain.License{
		ID:            "lic1",
		PublicationID: "pub1",
	}

	err := repo.Save(context.Background(), lic)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	found, err := repo.FindByID(context.Background(), "lic1")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected license")
	}

	err = repo.Delete(context.Background(), "lic1")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	found, err = repo.FindByID(context.Background(), "lic1")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}

	if found != nil {
		t.Fatal("expected deleted license")
	}
}

func TestPersistentPublicationRepositoryLoadsSavedData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "publications.json")

	repo, err := NewPersistentPublicationRepository(path)
	if err != nil {
		t.Fatalf("create persistent repo: %v", err)
	}

	err = repo.Save(context.Background(), &domain.Publication{
		ID:    "pub1",
		Title: "Persistent Book",
	})
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reloaded, err := NewPersistentPublicationRepository(path)
	if err != nil {
		t.Fatalf("reload persistent repo: %v", err)
	}

	found, err := reloaded.FindByID(context.Background(), "pub1")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if found == nil || found.Title != "Persistent Book" {
		t.Fatalf("unexpected publication: %#v", found)
	}
}

func TestPersistentLicenseRepositoryLoadsSavedData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "licenses.json")

	repo, err := NewPersistentLicenseRepository(path)
	if err != nil {
		t.Fatalf("create persistent repo: %v", err)
	}

	err = repo.Save(context.Background(), &domain.License{
		ID:            "lic1",
		PublicationID: "pub1",
	})
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reloaded, err := NewPersistentLicenseRepository(path)
	if err != nil {
		t.Fatalf("reload persistent repo: %v", err)
	}

	found, err := reloaded.FindByID(context.Background(), "lic1")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if found == nil || found.PublicationID != "pub1" {
		t.Fatalf("unexpected license: %#v", found)
	}
}

func TestLicenseRepositoryFindByPublication(t *testing.T) {
	repo := NewLicenseRepository()
	if err := repo.Save(context.Background(), &domain.License{ID: "lic1", PublicationID: "pub1"}); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if err := repo.Save(context.Background(), &domain.License{ID: "lic2", PublicationID: "pub2"}); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	publicationID := "pub1"
	items, err := repo.FindByPublication(context.Background(), &publicationID)
	if err != nil {
		t.Fatalf("find by publication failed: %v", err)
	}
	if len(items) != 1 || items[0].ID != "lic1" {
		t.Fatalf("unexpected filtered licenses: %+v", items)
	}

	all, err := repo.FindByPublication(context.Background(), nil)
	if err != nil {
		t.Fatalf("find all by publication failed: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected two licenses, got %d", len(all))
	}
}
