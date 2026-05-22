package tenant

import (
	"context"
	"testing"

	"github.com/amirhdev/ebook-lcp-server/internal/auth"
)

func TestIDFromContext(t *testing.T) {
	if got := IDFromContext(context.Background()); got != "default" {
		t.Fatalf("expected default tenant, got %q", got)
	}
	ctx := auth.WithClaims(context.Background(), &auth.Claims{TenantID: "tenant-a"})
	if got := IDFromContext(ctx); got != "tenant-a" {
		t.Fatalf("expected tenant from claims, got %q", got)
	}
}
