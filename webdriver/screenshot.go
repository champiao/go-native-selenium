package webdriver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Screenshot captura a tela atual do navegador e salva em um arquivo PNG.
// Screenshot captures the current browser screen and saves it to a PNG file.
func (c *Client) Screenshot(filename string) error {
	// Endpoint W3C para captura de tela
	// W3C endpoint for taking a screenshot
	endpoint := fmt.Sprintf("%s/session/%s/screenshot", c.BaseURL, c.SessionID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create screenshot request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("screenshot failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Value string `json:"value"` // A imagem vem como uma string Base64
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode screenshot response: %v", err)
	}

	// Decodifica a string Base64 para bytes puros de imagem
	// Decodes the Base64 string to pure image bytes
	imageBytes, err := base64.StdEncoding.DecodeString(result.Value)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Salva os bytes no disco com o nome fornecido
	// Saves the bytes to disk with the provided filename
	if err := os.WriteFile(filename, imageBytes, 0644); err != nil {
		return fmt.Errorf("failed to save screenshot to disk: %v", err)
	}

	return nil
}
