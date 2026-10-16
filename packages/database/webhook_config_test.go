package database

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type webhookRoundTrip func(*http.Request) (*http.Response, error)

func (f webhookRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestAdidasWebhookConfiguration(t *testing.T) {
	const key = "BONSAI_WEBHOOK_ADIDAS_LOG"
	calls := 0
	original := http.DefaultTransport
	http.DefaultTransport = webhookRoundTrip(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})
	defer func() { http.DefaultTransport = original }()

	t.Setenv(key, "")
	AdidasLog(Product{})
	if calls != 0 {
		t.Fatalf("missing configuration made %d requests", calls)
	}

	t.Setenv(key, "https://example.invalid/webhook")
	AdidasLog(Product{})
	if calls != 1 {
		t.Fatalf("expected one configured webhook request, got %d", calls)
	}
}
