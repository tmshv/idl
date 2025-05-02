package dl

import (
	"io"
	"net/http"
	"time"
)

type Downloader struct {
	Timeout time.Duration
	Retries int

	client *http.Client
}

func (dl *Downloader) Download(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := dl.client.Do(req)
	if err != nil {
		return nil, err
	}

	// var body []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}

func New(timeout time.Duration, retries int) *Downloader {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Downloader{
		Timeout: timeout,
		Retries: retries,
		client:  client,
	}
}
