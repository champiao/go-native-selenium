package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// Usaremos o endpoint de Milestones, que é muito mais estável para versões dinâmicas
// We will use the Milestones endpoint, which is much more stable for dynamic versions
const milestoneURL = "https://googlechromelabs.github.io/chrome-for-testing/latest-versions-per-milestone-with-downloads.json"

// Estrutura atualizada para mapear o JSON de milestones do Google
type milestoneResponse struct {
	Milestones map[string]struct {
		Version   string `json:"version"`
		Downloads struct {
			ChromeDriver []struct {
				Platform string `json:"platform"`
				URL      string `json:"url"`
			} `json:"chromedriver"`
		} `json:"downloads"`
	} `json:"milestones"`
}

// GetDriverURL busca o link de download baseado na versão principal (Major).
// GetDriverURL fetches the download link based on the Major version.
func GetDriverURL(targetVersion string) (string, error) {
	platform := getGooglePlatformName()
	if platform == "" {
		return "", errors.New("unsupported platform for chromedriver")
	}

	// Extrai apenas a versão principal (ex: "145" de "145.0.7632.116")
	// Extracts only the major version
	majorVersion := strings.Split(targetVersion, ".")[0]

	resp, err := http.Get(milestoneURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch chrome versions API: %v", err)
	}
	defer resp.Body.Close()

	var data milestoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %v", err)
	}

	// Busca os dados do Milestone específico (ex: "145")
	milestoneData, exists := data.Milestones[majorVersion]
	if !exists {
		return "", fmt.Errorf("chromedriver milestone %s not found", majorVersion)
	}

	// Encontra a URL para a plataforma correta
	for _, driverDownload := range milestoneData.Downloads.ChromeDriver {
		if driverDownload.Platform == platform {
			return driverDownload.URL, nil
		}
	}

	return "", fmt.Errorf("chromedriver for milestone %s not found for platform %s", majorVersion, platform)
}

// getGooglePlatformName converte a arquitetura do Go para o padrão do Google.
// getGooglePlatformName converts Go architecture to Google's format.
func getGooglePlatformName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "linux":
		return "linux64"
	case "darwin": // macOS
		if arch == "arm64" {
			return "mac-arm64"
		}
		return "mac-x64"
	case "windows":
		if arch == "amd64" {
			return "win64"
		}
		return "win32"
	}
	return ""
}
