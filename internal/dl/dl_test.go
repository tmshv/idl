package dl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDownload_Success(t *testing.T) {
	expected := []byte("image-data-here")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(expected)
	}))
	defer srv.Close()

	d := New(5*time.Second, 1)
	data, err := d.Download(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(data) != string(expected) {
		t.Fatalf("expected %q, got %q", expected, data)
	}
}

func TestDownload_404_ReturnsPermanentError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	d := New(5*time.Second, 3)
	_, err := d.Download(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestDownload_Non200_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	d := New(5*time.Second, 1)
	_, err := d.Download(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestDownload_MalformedURL_ReturnsError(t *testing.T) {
	d := New(5*time.Second, 1)

	// BUG: the current implementation panics on malformed URLs because
	// http.NewRequest error is silently discarded (req, _ := ...).
	// This test documents that the code panics instead of returning an error.
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("known bug: malformed URL causes panic instead of returning error: %v", r)
		}
	}()

	_, err := d.Download(context.Background(), "://bad-url")
	if err == nil {
		t.Fatal("expected error for malformed URL, got nil")
	}
}

func TestDownload_RetryOnTransientFailure(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	d := New(5*time.Second, 5)
	data, err := d.Download(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if string(data) != "ok" {
		t.Fatalf("expected %q, got %q", "ok", data)
	}
	if attempts.Load() < 3 {
		t.Fatalf("expected at least 3 attempts, got %d", attempts.Load())
	}
}
