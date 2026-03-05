package webdriver

// Capabilities representa a configuração desejada para a sessão do navegador.
// Capabilities represents the desired configuration for the browser session.
type Capabilities struct {
	BrowserName   string
	ChromeOptions *ChromeOptions
}

// ChromeOptions mapeia as configurações específicas da engine do Google Chrome.
// ChromeOptions maps the specific configurations for the Google Chrome engine.
type ChromeOptions struct {
	Args []string `json:"args,omitempty"`
}

// DefaultCapabilities retorna uma configuração padrão limpa para o Chrome.
// DefaultCapabilities returns a clean standard configuration for Chrome.
func DefaultCapabilities() *Capabilities {
	return &Capabilities{
		BrowserName: "chrome",
		ChromeOptions: &ChromeOptions{
			Args: []string{},
		},
	}
}

// AddArgument adiciona uma flag de linha de comando ao Chrome (ex: "--headless").
// AddArgument adds a command-line flag to Chrome (e.g., "--headless").
func (c *Capabilities) AddArgument(arg string) {
	if c.ChromeOptions == nil {
		c.ChromeOptions = &ChromeOptions{}
	}
	c.ChromeOptions.Args = append(c.ChromeOptions.Args, arg)
}
