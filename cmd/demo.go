package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/champiao/go-native-selenium/driver"
	"github.com/champiao/go-native-selenium/webdriver"
)

func main() {
	fmt.Println("1. Detectando versão do Chrome...")
	version, err := driver.GetChromeVersion()
	if err != nil {
		log.Fatalf("Erro ao detectar versão: %v", err)
	}
	fmt.Printf("   Versão encontrada: %s\n", version)

	fmt.Println("2. Buscando URL de download...")
	majorVersion := strings.Split(version, ".")[0]
	url, err := driver.GetDriverURL(version)
	if err != nil {
		log.Fatalf("Erro ao buscar URL: %v", err)
	}
	fmt.Printf("   URL encontrada: %s\n", url)

	fmt.Println("3. Baixando e extraindo o driver...")
	driverPath, err := driver.DownloadDriver(url, majorVersion)
	if err != nil {
		log.Fatalf("Erro ao gerenciar o download: %v", err)
	}
	fmt.Printf("   Sucesso! Driver pronto para uso em: %s\n", driverPath)
	fmt.Println("4. Iniciando o ChromeDriver em background...")
	service, err := driver.StartService(driverPath)
	if err != nil {
		log.Fatalf("Erro ao iniciar o serviço: %v", err)
	}

	fmt.Printf("   Serviço rodando na porta: %d\n", service.Port)
	fmt.Printf("   Endpoint base da API: %s\n", service.URL)

	fmt.Println("5. Abrindo o navegador (Criando Sessão W3C)...")
	caps := webdriver.DefaultCapabilities()
	caps.AddArgument("--incognito")
	caps.AddArgument("--disable-gpu")
	client, err := webdriver.NewSession(service.URL, caps)
	if err != nil {
		log.Fatalf("Erro ao abrir navegador: %v", err)
	}
	if err := client.SetImplicitWait(10 * time.Second); err != nil {
		log.Printf("Aviso: Falha ao configurar wait: %v\n", err)
	}
	fmt.Println("   Maximizando a janela...")
	if err := client.MaximizeWindow(); err != nil {
		log.Printf("Aviso: Falha ao maximizar a janela: %v\n", err)
	}
	fmt.Println("   Navegando para o seu projeto...")
	if err := client.Navigate("https://github.com/champiao/go-native-selenium"); err != nil {
		log.Fatalf("Erro ao navegar: %v", err)
	}
	client.Screenshot("github_of_project.png")
	fmt.Println("   Evidência salva! Verifique o arquivo evidencia-dashboard-logado.png.")
	fmt.Println("   Página carregada!")
	fmt.Println("🚀 Teste finalizado com sucesso!")

	// ---------------------------------------------------------
	// MANTER SESSÃO ATIVA
	// ---------------------------------------------------------
	fmt.Println("\n[PAUSA] O navegador permanecerá aberto.")
	fmt.Println("👉 Pressione ENTER no terminal para encerrar o teste...")
	fmt.Scanln()

	fmt.Println("6. Fechando o navegador...")
	if err := client.Quit(); err != nil {
		log.Printf("Erro ao fechar o navegador: %v\n", err)
	}

	fmt.Println("7. Encerrando o serviço do ChromeDriver...")
	if err := service.Stop(); err != nil {
		log.Fatalf("Erro ao parar o serviço: %v", err)
	}
	fmt.Println("   Tudo encerrado com sucesso!")
}
