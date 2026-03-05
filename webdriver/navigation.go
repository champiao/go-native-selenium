package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Navigate instrui o navegador a carregar uma URL específica.
// Navigate instructs the browser to load a specific URL.
func (c *Client) Navigate(url string) error {
	payload := map[string]string{
		"url": url,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal url payload: %v", err)
	}

	// Endpoint W3C para navegação
	// W3C endpoint for navigation
	endpoint := fmt.Sprintf("%s/session/%s/url", c.BaseURL, c.SessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to navigate to url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("navigation failed with status: %d", resp.StatusCode)
	}

	return nil
}
