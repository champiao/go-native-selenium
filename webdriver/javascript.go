package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ExecuteScript executa um script JavaScript de forma síncrona no contexto da página atual.
// ExecuteScript synchronously executes a JavaScript script in the context of the current page.
func (c *Client) ExecuteScript(script string, args ...interface{}) (interface{}, error) {
	// Se args for nil, o JSON Marshal pode enviar null em vez de um array vazio [],
	// o que causaria um erro no ChromeDriver. Garantimos que seja um slice vazio.
	if args == nil {
		args = make([]interface{}, 0)
	}

	payload := map[string]interface{}{
		"script": script,
		"args":   args,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal script payload: %v", err)
	}

	// Endpoint W3C para execução de script síncrono
	// W3C endpoint for synchronous script execution
	endpoint := fmt.Sprintf("%s/session/%s/execute/sync", c.BaseURL, c.SessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("script execution failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	// O script pode retornar qualquer coisa (um número, uma string, um booleano ou nada).
	// Por isso usamos interface{}.
	var result struct {
		Value interface{} `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode script response: %v", err)
	}

	return result.Value, nil
}
