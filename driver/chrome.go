package driver

import (
	"errors"
	"os/exec"
	"regexp"
	"runtime"
)

// GetChromeVersion detecta a versão do Google Chrome instalada no sistema.
// GetChromeVersion detects the installed Google Chrome version on the system.
func GetChromeVersion() (string, error) {
	var out []byte
	var err error

	switch runtime.GOOS {
	case "windows":
		// Consulta o registro do Windows / Queries the Windows registry
		out, err = exec.Command("reg", "query", `HKEY_CURRENT_USER\Software\Google\Chrome\BLBeacon`, "/v", "version").Output()
	case "darwin":
		// Caminho padrão no macOS / Default path on macOS
		out, err = exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--version").Output()
	case "linux":
		// Tenta os binários mais comuns em distribuições Linux / Tries common binaries on Linux distros
		out, err = exec.Command("google-chrome", "--version").Output()
		if err != nil {
			out, err = exec.Command("google-chrome-stable", "--version").Output()
		}
	default:
		return "", errors.New("unsupported operating system")
	}

	if err != nil {
		return "", errors.New("could not find chrome installation")
	}

	return extractVersion(string(out))
}

// extractVersion isola apenas os números da versão (ex: 122.0.6261.94).
// extractVersion isolates only the version numbers.
func extractVersion(output string) (string, error) {
	re := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	match := re.FindString(output)
	if match == "" {
		return "", errors.New("version format not found in output: " + output)
	}
	return match, nil
}
