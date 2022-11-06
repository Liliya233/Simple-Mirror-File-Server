package utils

import (
	"bytes"
	"io"
	"net/http"
)

func DownloadSmallContent(sourceUrl string) []byte {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Buffer
	contents := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(contents, resp.Body); err == nil {
		return contents.Bytes()
	} else {
		panic(err)
	}
}
