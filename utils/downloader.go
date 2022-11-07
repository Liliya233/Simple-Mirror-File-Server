package utils

import (
	"bytes"
	"io"
	"net/http"
)

func DownloadSmallContent(sourceUrl string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Buffer
	contents := bytes.NewBuffer([]byte{})
	_, err = io.Copy(contents, resp.Body)

	// Return bytes or err
	return contents.Bytes(), err
}
