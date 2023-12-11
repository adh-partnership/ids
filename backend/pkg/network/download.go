package network

import (
	"io"
	"net/http"

	"github.com/adh-partnership/api/pkg/logger"
)

func DownloadFile(url string) ([]byte, error) {
	logger.Logger.WithField("component", "network").Infof("downloading %s", url)
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
