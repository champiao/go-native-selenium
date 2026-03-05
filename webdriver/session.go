package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client representa a conexão ativa com o navegador através do ChromeDriver.
// Client represents the active connection to the browser via ChromeDriver.
type Client struct {
	BaseURL   string
	SessionID string
}

// w3cResponse mapeia o padrão de resposta do protocolo W3C.
// w3cResponse maps the W3C protocol response pattern.
type w3cResponse struct {
	Value struct {
		SessionID string `json:"sessionId"`
	} `json:"value"`
}

// NewSession envia o comando para o ChromeDriver abrir a janela do navegador.
// NewSession sends the command to ChromeDriver to open the browser window.
func NewSession(driverURL string, caps *Capabilities) (*Client, error) {
	if caps == nil {
		caps = DefaultCapabilities()
	}

	// Constrói o bloco "alwaysMatch" dinamicamente
	// Builds the "alwaysMatch" block dynamically
	alwaysMatch := map[string]interface{}{
		"browserName": caps.BrowserName,
	}

	// Injeta as opções específicas do Chrome se existirem
	// Injects specific Chrome options if they exist
	if caps.ChromeOptions != nil && len(caps.ChromeOptions.Args) > 0 {
		alwaysMatch["goog:chromeOptions"] = map[string]interface{}{
			"args": caps.ChromeOptions.Args,
		}
	}

	payload := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": alwaysMatch,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capabilities: %v", err)
	}

	sessionURL := fmt.Sprintf("%s/session", driverURL)
	resp, err := http.Post(sessionURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create session, status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var result w3cResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode session response: %v", err)
	}

	return &Client{
		BaseURL:   driverURL,
		SessionID: result.Value.SessionID,
	}, nil
}

// Quit envia o comando para fechar o navegador e destruir a sessão.
// Quit sends the command to close the browser and destroy the session.
func (c *Client) Quit() error {
	quitURL := fmt.Sprintf("%s/session/%s", c.BaseURL, c.SessionID)
	req, err := http.NewRequest(http.MethodDelete, quitURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to quit session, status: %d", resp.StatusCode)
	}
	return nil
}
