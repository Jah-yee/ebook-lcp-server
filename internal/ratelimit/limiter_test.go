package ratelimit

import (
	"testing"
	"time"
)

func TestAllowRespectsLimit(t *testing.T) {
	limiter := New(2, time.Minute)
	if !limiter.Allow("user-1") {
		t.Fatal("expected first request to pass")
	}
	if !limiter.Allow("user-1") {
		t.Fatal("expected second request to pass")
	}
	if limiter.Allow("user-1") {
		t.Fatal("expected third request to be blocked")
	}
}

func TestAllowResetsWindow(t *testing.T) {
	limiter := New(1, time.Minute)
	if !limiter.Allow("user-1") {
		t.Fatal("expected initial request to pass")
	}
	limiter.resetTime = time.Now().Add(-time.Second)
	if !limiter.Allow("user-1") {
		t.Fatal("expected request after reset window to pass")
	}
}

func TestAllowWithDisabledLimiterAlwaysPasses(t *testing.T) {
	var limiter *Limiter
	if !limiter.Allow("user-1") {
		t.Fatal("expected nil limiter to allow")
	}
	if !New(0, time.Minute).Allow("user-1") {
		t.Fatal("expected zero-limit limiter to allow")
	}
}
