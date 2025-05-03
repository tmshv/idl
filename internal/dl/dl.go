package dl

import (
	"errors"
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
	var e error
	for range dl.Retries {
		body, err := dl.dl(url)
		if err != nil {
			e = err
			continue
		}
		return body, nil
	}
	return nil, e
}

func (dl *Downloader) dl(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := dl.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not 200 OK")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

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
