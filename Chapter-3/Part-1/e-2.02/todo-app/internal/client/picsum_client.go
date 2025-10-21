package client

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PicsumClient struct {
	HTTPClient
	ImageDir string
}

func NewPicsumClient(baseUrl string, imageDir string, timeout time.Duration) *PicsumClient {
	return &PicsumClient{
		HTTPClient: HTTPClient{
			client: &http.Client{
				Timeout: timeout,
			},
			baseURL: baseUrl,
		},
		ImageDir: imageDir,
	}
}

func (pc *PicsumClient) GetClientBaseURL() string {
	return pc.baseURL
}

func (pc *PicsumClient) FetchAndSaveImage() error {
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(pc.ImageDir, 0o755); err != nil {
		return err
	}

	resp, err := http.Get(pc.baseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(pc.ImageDir, "background.jpg")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
