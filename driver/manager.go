package driver

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"time"
)

// Service representa a instância do ChromeDriver rodando em background.
// Service represents the ChromeDriver instance running in the background.
type Service struct {
	cmd  *exec.Cmd
	Port int
	URL  string
}

// StartService inicia o ChromeDriver em uma porta disponível no sistema.
// StartService starts the ChromeDriver on an available system port.
func StartService(binPath string) (*Service, error) {
	// 1. Encontra uma porta TCP livre / Finds a free TCP port
	port, err := getFreePort()
	if err != nil {
		return nil, fmt.Errorf("could not find a free port: %v", err)
	}

	// 2. Prepara o comando / Prepares the command
	// Ex: /home/.../chromedriver --port=9515
	cmd := exec.Command(binPath, fmt.Sprintf("--port=%d", port))

	// Opcional: Ocultar a janela do terminal no Windows (não afeta Linux/macOS)
	// hideWindow(cmd) // Implementaremos isso depois se necessário

	// 3. Inicia o processo de forma assíncrona / Starts the process asynchronously
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start chromedriver: %v", err)
	}

	// 4. Aguarda o servidor subir / Waits for the server to boot
	// (O ideal em produção é fazer um health-check HTTP no endpoint /status)
	time.Sleep(1 * time.Second)

	return &Service{
		cmd:  cmd,
		Port: port,
		URL:  fmt.Sprintf("http://localhost:%d", port),
	}, nil
}

// Stop encerra o processo do ChromeDriver de forma segura.
// Stop safely terminates the ChromeDriver process.
func (s *Service) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		// Envia o sinal de Kill para o processo do SO
		// Sends the Kill signal to the OS process
		return s.cmd.Process.Kill()
	}
	return errors.New("service is not running")
}

// getFreePort pede ao sistema operacional uma porta TCP aleatória livre.
// getFreePort asks the operating system for a random free TCP port.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
