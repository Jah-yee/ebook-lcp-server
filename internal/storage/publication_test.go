package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseS3URI(t *testing.T) {
	bucket, key, err := parseS3URI("s3://books/publications/book.epub")
	if err != nil {
		t.Fatalf("parseS3URI failed: %v", err)
	}
	if bucket != "books" || key != "publications/book.epub" {
		t.Fatalf("unexpected result: bucket=%s key=%s", bucket, key)
	}
}

func TestParseS3URIRejectsInvalidValue(t *testing.T) {
	if _, _, err := parseS3URI("https://example.test/book.epub"); err == nil {
		t.Fatal("expected invalid scheme to fail")
	}
}

func TestFilesystemStorageDoesNotSignURLs(t *testing.T) {
	url, ok, err := NewFilesystemPublicationStorage().SignedURL(context.Background(), "/tmp/book.epub", time.Minute)
	if err != nil {
		t.Fatalf("SignedURL failed: %v", err)
	}
	if ok || url != "" {
		t.Fatalf("expected no signed URL, got ok=%v url=%q", ok, url)
	}
}

func TestFilesystemStorageStoreAndOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "book.epub")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	store := NewFilesystemPublicationStorage()
	uri, err := store.StoreEncrypted(context.Background(), path, "pub-1")
	if err != nil {
		t.Fatalf("StoreEncrypted failed: %v", err)
	}
	if uri != path {
		t.Fatalf("expected local path URI, got %q", uri)
	}
	reader, err := store.OpenEncrypted(context.Background(), uri)
	if err != nil {
		t.Fatalf("OpenEncrypted failed: %v", err)
	}
	defer func() { _ = reader.Close() }()
	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected file body %q", string(body))
	}
	if err := store.Ready(context.Background()); err != nil {
		t.Fatalf("Ready failed: %v", err)
	}
}
