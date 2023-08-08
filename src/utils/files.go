package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func AddItem(f Image) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", f.URI, nil)
	filename, err := generateFilename(f.URI, req.URL.Path)
	if err != nil {
		return "", fmt.Errorf("Unable to create filename %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Unable to fetch image %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to parse http response body %v", err)
	}

	localPath := os.Getenv("WISH_FILE_PATH")
	err = os.WriteFile(localPath+"/"+filename, body, 0666)
	if err != nil {
		return "", fmt.Errorf("Unable to save image %v", err)
	}

	return filename, nil
}

func generateFilename(uri string, path string) (string, error) {
	hash := sha256.Sum256([]byte(uri))

	split := strings.Split(path, ".")

	ext := split[len(split)-1]
	for _, e := range []string{"jpg", "jpeg", "png"} {
		if e == ext {
			return hex.EncodeToString(hash[:]) + "." + e, nil
		}
	}
	return "", fmt.Errorf("Unable to parse filename of image, check if it is a valid jpg, jpeg or png image.")
}
