package webdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Constante padrão do W3C para identificar o ID do elemento retornado
// Standard W3C constant to identify the returned element ID
const w3cElementKey = "element-6066-11e4-a52e-4f735466cecf"

// Estratégias de busca (By) similares ao Python
// Search strategies (By) similar to Python
const (
	ByID          = "id"
	ByName        = "name"
	ByClassName   = "class name"
	ByCSSSelector = "css selector"
	ByXPath       = "xpath"
	ByTagName     = "tag name"
)

// Element representa um nó DOM na página.
// Element represents a DOM node on the page.
type Element struct {
	ID     string
	client *Client
}

// w3cElementResponse mapeia a resposta de busca de elemento.
// w3cElementResponse maps the element search response.
type w3cElementResponse struct {
	Value map[string]string `json:"value"`
}

// FindElement busca um elemento, traduzindo estratégias amigáveis para o padrão W3C.
// FindElement searches for an element, translating friendly strategies to the W3C standard.
func (c *Client) FindElement(using, value string) (*Element, error) {

	// A "Mágica" do Python: Traduzimos os seletores amigáveis para CSS Selectors reais
	// The Python "Magic": We translate friendly selectors to actual CSS Selectors
	switch using {
	case ByID:
		using = ByCSSSelector
		value = fmt.Sprintf("#%s", value)
	case ByName:
		using = ByCSSSelector
		value = fmt.Sprintf("[name='%s']", value)
	case ByClassName:
		using = ByCSSSelector
		value = fmt.Sprintf(".%s", value)
	}

	payload := map[string]string{
		"using": using,
		"value": value,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal element payload: %v", err)
	}

	// Endpoint W3C para buscar um elemento
	// W3C endpoint to find an element
	endpoint := fmt.Sprintf("%s/session/%s/element", c.BaseURL, c.SessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to find element: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("element not found (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result w3cElementResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode element response: %v", err)
	}

	// O protocolo W3C retorna o ID do elemento em uma chave UUID específica
	// The W3C protocol returns the element ID in a specific UUID key
	elementID, ok := result.Value[w3cElementKey]
	if !ok {
		// Fallback para versões muito antigas do protocolo JSON Wire
		// Fallback for very old versions of the JSON Wire protocol
		elementID = result.Value["ELEMENT"]
	}

	return &Element{
		ID:     elementID,
		client: c, // Guardamos a referência do cliente para usar nos próximos cliques
	}, nil
}

// Click simula um clique do mouse no elemento.
// Click simulates a mouse click on the element.
func (e *Element) Click() error {
	// W3C endpoint: POST /session/{session id}/element/{element id}/click
	endpoint := fmt.Sprintf("%s/session/%s/element/%s/click", e.client.BaseURL, e.client.SessionID, e.ID)

	// O protocolo W3C exige um objeto JSON vazio no corpo do POST para o clique
	// The W3C protocol requires an empty JSON object in the POST body for clicking
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return fmt.Errorf("failed to click element: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("click failed with status: %d", resp.StatusCode)
	}
	return nil
}

// SendKeys simula a digitação de texto em um campo de input.
// SendKeys simulates typing text into an input field.
func (e *Element) SendKeys(text string) error {
	// W3C endpoint: POST /session/{session id}/element/{element id}/value
	endpoint := fmt.Sprintf("%s/session/%s/element/%s/value", e.client.BaseURL, e.client.SessionID, e.ID)

	payload := map[string]string{
		"text": text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal keys payload: %v", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send keys to element: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send keys failed with status: %d", resp.StatusCode)
	}
	return nil
}

// GetText retorna o texto visível interno de um elemento.
// GetText returns the visible inner text of an element.
func (e *Element) GetText() (string, error) {
	// W3C endpoint: GET /session/{session id}/element/{element id}/text
	endpoint := fmt.Sprintf("%s/session/%s/element/%s/text", e.client.BaseURL, e.client.SessionID, e.ID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create get text request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get text from element: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get text failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode text response: %v", err)
	}

	return result.Value, nil
}

// GetAttribute retorna o valor de um atributo HTML específico do elemento.
// GetAttribute returns the value of a specific HTML attribute of the element.
func (e *Element) GetAttribute(name string) (string, error) {
	// W3C endpoint: GET /session/{session id}/element/{element id}/attribute/{name}
	endpoint := fmt.Sprintf("%s/session/%s/element/%s/attribute/%s", e.client.BaseURL, e.client.SessionID, e.ID, name)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create get attribute request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get attribute from element: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get attribute failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode attribute response: %v", err)
	}

	return result.Value, nil
}

// w3cElementsResponse mapeia a resposta para múltiplos elementos.
// w3cElementsResponse maps the response for multiple elements.
type w3cElementsResponse struct {
	Value []map[string]string `json:"value"`
}

// FindElements busca múltiplos elementos na página e retorna uma fatia (slice) de Elements.
// FindElements searches for multiple elements on the page and returns a slice of Elements.
func (c *Client) FindElements(using, value string) ([]*Element, error) {
	// A mesma "Mágica" de tradução do Python
	// The same Python translation "Magic"
	switch using {
	case ByID:
		using = ByCSSSelector
		value = fmt.Sprintf("#%s", value)
	case ByName:
		using = ByCSSSelector
		value = fmt.Sprintf("[name='%s']", value)
	case ByClassName:
		using = ByCSSSelector
		value = fmt.Sprintf(".%s", value)
	}

	payload := map[string]string{
		"using": using,
		"value": value,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal elements payload: %v", err)
	}

	// Endpoint W3C no plural: /elements
	// W3C plural endpoint: /elements
	endpoint := fmt.Sprintf("%s/session/%s/elements", c.BaseURL, c.SessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to find elements: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("elements not found (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result w3cElementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode elements response: %v", err)
	}

	// Constrói o slice de elementos Go a partir dos IDs retornados
	// Builds the Go elements slice from the returned IDs
	var elements []*Element
	for _, item := range result.Value {
		elementID, ok := item[w3cElementKey]
		if !ok {
			elementID = item["ELEMENT"] // Fallback
		}
		elements = append(elements, &Element{
			ID:     elementID,
			client: c,
		})
	}

	return elements, nil
}
