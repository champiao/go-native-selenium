package webdriver

import (
	"fmt"
	"time"
)

// Wait define as configurações para aguardar uma condição na tela.
// Wait defines the configuration for waiting for a condition on the screen.
type Wait struct {
	client       *Client
	Timeout      time.Duration
	PollInterval time.Duration
}

// NewWait cria uma nova instância de espera explícita.
// NewWait creates a new explicit wait instance.
func (c *Client) NewWait(timeout time.Duration) *Wait {
	return &Wait{
		client:       c,
		Timeout:      timeout,
		PollInterval: 500 * time.Millisecond, // Intervalo padrão entre as tentativas
	}
}

// ExpectedCondition é uma assinatura de função que retorna um elemento e um erro se falhar.
// ExpectedCondition is a function signature that returns an element and an error if it fails.
type ExpectedCondition func() (*Element, error)

// UntilElementLocated retorna uma condição que verifica se o elemento já existe no DOM.
// UntilElementLocated returns a condition that checks if the element exists in the DOM.
func (c *Client) UntilElementLocated(using, value string) ExpectedCondition {
	return func() (*Element, error) {
		return c.FindElement(using, value)
	}
}

// Until executa a condição repetidamente até que ela tenha sucesso ou o tempo acabe.
// Until executes the condition repeatedly until it succeeds or times out.
func (w *Wait) Until(condition ExpectedCondition) (*Element, error) {
	endTime := time.Now().Add(w.Timeout)

	for time.Now().Before(endTime) {
		el, err := condition()
		if err == nil && el != nil {
			return el, nil // Sucesso! Elemento encontrado.
		}

		// Aguarda o intervalo antes de tentar novamente
		time.Sleep(w.PollInterval)
	}

	return nil, fmt.Errorf("timeout of %v exceeded waiting for condition", w.Timeout)
}
