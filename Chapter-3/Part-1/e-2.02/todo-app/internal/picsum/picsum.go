package picsum

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const ImageURL = "https://picsum.photos/1200"

// FetchAndSaveImage fetches an image from the specified URL and saves it to the given directory.
func FetchAndSaveImage(dir string) error {
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	resp, err := http.Get(ImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(dir, "background.jpg")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
