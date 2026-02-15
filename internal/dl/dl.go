package dl

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
)

type Downloader struct {
	Timeout time.Duration
	retries uint

	client *http.Client
}

func (dl *Downloader) Download(ctx context.Context, url string) ([]byte, error) {
	operation := func() ([]byte, error) {
		return dl.dl(url)
	}
	return backoff.Retry(ctx, operation,
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(dl.retries),
	)
}

func (dl *Downloader) dl(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, backoff.Permanent(err)
	}
	res, err := dl.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	// In case on non-retriable error, return Permanent error to stop retrying.
	// For this HTTP example, client errors are non-retriable.
	if res.StatusCode == 404 {
		return nil, backoff.Permanent(errors.New("not found"))
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not 200 OK")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func New(timeout time.Duration, retries uint) *Downloader {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Downloader{
		Timeout: timeout,
		retries: retries,
		client:  client,
	}
}
