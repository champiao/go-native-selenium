package driver

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// DownloadDriver baixa e extrai o ChromeDriver para o diretório de cache do sistema.
// DownloadDriver downloads and extracts ChromeDriver to the system cache directory.
func DownloadDriver(url, version string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not get user cache dir: %v", err)
	}

	// Cria o caminho: ~/.cache/go-selenium-native/chromedriver/<version>/
	// Creates the path: ~/.cache/go-selenium-native/chromedriver/<version>/
	targetDir := filepath.Join(cacheDir, "go-selenium-native", "chromedriver", version)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Define o nome do executável esperado
	// Defines the expected executable name
	binaryName := "chromedriver"
	if runtime.GOOS == "windows" {
		binaryName = "chromedriver.exe"
	}

	targetPath := filepath.Join(targetDir, binaryName)

	// Verifica se o binário já existe (Cache Hit)
	// Checks if the binary already exists (Cache Hit)
	if _, err := os.Stat(targetPath); err == nil {
		return targetPath, nil // Já baixado! / Already downloaded!
	}

	// Faz o download do arquivo ZIP
	// Downloads the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download driver: %v", err)
	}
	defer resp.Body.Close()

	// Lê o ZIP inteiro para a memória (é pequeno, ~8MB)
	// Reads the entire ZIP into memory (it's small, ~8MB)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", fmt.Errorf("failed to read zip archive: %v", err)
	}

	// Procura o executável dentro do ZIP
	// Searches for the executable inside the ZIP
	for _, file := range zipReader.File {
		if filepath.Base(file.Name) == binaryName {
			if err := extractFile(file, targetPath); err != nil {
				return "", err
			}
			break
		}
	}

	return targetPath, nil
}

// extractFile extrai um único arquivo do ZIP e define permissões de execução.
// extractFile extracts a single file from the ZIP and sets execution permissions.
func extractFile(file *zip.File, destPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, rc)
	return err
}
