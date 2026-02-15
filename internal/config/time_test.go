package config

import (
	"testing"
	"time"
)

func TestParseTimeout_EmptyString(t *testing.T) {
	d, err := ParseTimeout("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 0 {
		t.Fatalf("expected 0, got %v", d)
	}
}

func TestParseTimeout_Seconds(t *testing.T) {
	d, err := ParseTimeout("10s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 10*time.Second {
		t.Fatalf("expected 10s, got %v", d)
	}
}

func TestParseTimeout_Minutes(t *testing.T) {
	d, err := ParseTimeout("1m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 1*time.Minute {
		t.Fatalf("expected 1m, got %v", d)
	}
}

func TestParseTimeout_BareNumber(t *testing.T) {
	d, err := ParseTimeout("10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 10*time.Second {
		t.Fatalf("expected 10s, got %v", d)
	}
}

func TestParseTimeout_InvalidInput(t *testing.T) {
	_, err := ParseTimeout("abc")
	if err == nil {
		t.Fatal("expected error for invalid input, got nil")
	}
}
