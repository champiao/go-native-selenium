package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MaximizeWindow expande a janela do navegador para ocupar toda a tela.
// MaximizeWindow expands the browser window to fill the entire screen.
func (c *Client) MaximizeWindow() error {
	// Endpoint W3C para maximizar a janela
	// W3C endpoint to maximize the window
	endpoint := fmt.Sprintf("%s/session/%s/window/maximize", c.BaseURL, c.SessionID)

	// O protocolo W3C exige um payload JSON vazio "{}"
	// The W3C protocol requires an empty JSON payload "{}"
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return fmt.Errorf("failed to maximize window: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("maximize window failed with status: %d", resp.StatusCode)
	}
	return nil
}

// SetWindowSize define uma resolução específica para a janela (largura e altura).
// SetWindowSize sets a specific resolution for the window (width and height).
func (c *Client) SetWindowSize(width, height int) error {
	// Endpoint W3C para alterar as dimensões da janela (Rect)
	// W3C endpoint to change window dimensions (Rect)
	endpoint := fmt.Sprintf("%s/session/%s/window/rect", c.BaseURL, c.SessionID)

	payload := map[string]int{
		"width":  width,
		"height": height,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal window size payload: %v", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to set window size: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set window size failed with status: %d", resp.StatusCode)
	}
	return nil
}
