package requestmeta

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIDFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestIDKey, "req-1")
	if got := IDFromContext(ctx); got != "req-1" {
		t.Fatalf("expected request id, got %q", got)
	}
	if got := IDFromContext(context.Background()); got != "" {
		t.Fatalf("expected empty request id, got %q", got)
	}
}

func TestMiddlewareUsesIncomingRequestIDAndWritesHeader(t *testing.T) {
	var body strings.Builder
	logger := slog.New(slog.NewTextHandler(&body, nil))
	handler := Middleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := IDFromContext(r.Context()); got != "external-id" {
			t.Fatalf("expected request id in context, got %q", got)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "ok")
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-ID", "external-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); got != "external-id" {
		t.Fatalf("expected response request id header, got %q", got)
	}
	if !strings.Contains(body.String(), "request_id=external-id") {
		t.Fatalf("expected request id in logs, got %q", body.String())
	}
}

func TestMiddlewareGeneratesRequestIDWhenMissing(t *testing.T) {
	handler := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := IDFromContext(r.Context()); len(got) != 32 {
			t.Fatalf("expected generated request id, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); len(got) != 32 {
		t.Fatalf("expected generated response request id, got %q", got)
	}
}

func TestStatusRecorderWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	w := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}
	w.WriteHeader(http.StatusAccepted)
	if w.status != http.StatusAccepted {
		t.Fatalf("expected recorded status %d, got %d", http.StatusAccepted, w.status)
	}
}

func TestNewIDReturnsHexString(t *testing.T) {
	got := newID()
	if len(got) != 32 {
		t.Fatalf("expected 32-char id, got %q", got)
	}
}
