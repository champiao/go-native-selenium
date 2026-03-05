package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SetImplicitWait define o tempo máximo que o WebDriver deve tentar buscar um elemento.
// SetImplicitWait sets the maximum time the WebDriver should try to find an element.
func (c *Client) SetImplicitWait(timeout time.Duration) error {
	// Endpoint W3C para configurar os timeouts da sessão
	// W3C endpoint to configure session timeouts
	endpoint := fmt.Sprintf("%s/session/%s/timeouts", c.BaseURL, c.SessionID)

	// O W3C espera o tempo em milissegundos
	// W3C expects the time in milliseconds
	payload := map[string]int64{
		"implicit": timeout.Milliseconds(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal timeouts payload: %v", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to set implicit wait: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set implicit wait failed with status: %d", resp.StatusCode)
	}

	return nil
}
