package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Cookie representa um HTTP Cookie padrão do W3C.
// Cookie represents a standard W3C HTTP Cookie.
type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HttpOnly bool   `json:"httpOnly,omitempty"`
	Expiry   int64  `json:"expiry,omitempty"`
}

// AddCookie injeta um cookie no navegador atual.
// IMPORTANTE: O navegador já deve estar no domínio do cookie antes de injetá-lo.
// AddCookie injects a cookie into the current browser.
// IMPORTANT: The browser must already be on the cookie's domain before injecting it.
func (c *Client) AddCookie(cookie *Cookie) error {
	payload := map[string]*Cookie{
		"cookie": cookie,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal cookie payload: %v", err)
	}

	// Endpoint W3C para adicionar um cookie
	// W3C endpoint to add a cookie
	endpoint := fmt.Sprintf("%s/session/%s/cookie", c.BaseURL, c.SessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to add cookie: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add cookie failed with status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GetCookies retorna todos os cookies associados à página atual.
// GetCookies returns all cookies associated with the current page.
func (c *Client) GetCookies() ([]*Cookie, error) {
	endpoint := fmt.Sprintf("%s/session/%s/cookie", c.BaseURL, c.SessionID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get cookies request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cookies failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Value []*Cookie `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode cookies response: %v", err)
	}

	return result.Value, nil
}
