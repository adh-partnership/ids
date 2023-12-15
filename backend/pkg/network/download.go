package network

import (
	"io"
	"net/http"
)

func Request(method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("User-Agent", "ADH Partnership IDS")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}
