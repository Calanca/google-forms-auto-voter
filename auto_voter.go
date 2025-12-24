package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// Entry ID para a primeira porta
	entryID   = "entry.877086558"
	voteValue = "Amei esta!"
)

// Lê o link de votação do arquivo config.txt
func readFormURLFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignora linhas vazias e comentários
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			// Converte viewform para formResponse se necessário
			line = strings.Replace(line, "/viewform", "/formResponse", 1)
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %v", err)
	}

	return "", fmt.Errorf("nenhuma URL válida encontrada no arquivo")
}

func submitVote(voteNumber int, formURL string) error {
	// Cria os dados do formulário
	data := url.Values{}
	data.Set(entryID, voteValue)

	// Cria a requisição POST
	req, err := http.NewRequest("POST", formURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}

	// Define os headers necessários
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	// Converte formResponse para viewform no Referer
	refererURL := strings.Replace(formURL, "/formResponse", "/viewform", 1)
	req.Header.Set("Referer", refererURL)

	// Cria o cliente HTTP que SEGUE redirects
	var finalURL string
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Captura a URL final após o redirect
			finalURL = req.URL.String()
			return nil
		},
	}

	// Envia a requisição
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar voto: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %v", err)
	}

	bodyStr := string(body)

	// Verifica se a resposta foi bem-sucedida
	if resp.StatusCode == http.StatusOK {
		// Verifica se chegou na página de confirmação
		if strings.Contains(finalURL, "formResponse") || strings.Contains(resp.Request.URL.String(), "formResponse") {
			// Verifica se contém a mensagem exata de sucesso do formulário
			if strings.Contains(bodyStr, "Obrigada por escolher a portinha mais fashion do HDSA") ||
				strings.Contains(bodyStr, "Edite a sua resposta") ||
				strings.Contains(bodyStr, "Enviar outra resposta") {
				fmt.Printf("✓ Voto #%d CONFIRMADO! 🎄 Mensagem: 'Obrigada por escolher a portinha mais fashion do HDSA ;D'\n", voteNumber)
				return nil
			}
		}

		// Se chegou aqui, algo pode estar errado
		fmt.Printf("⚠ Voto #%d: Status OK mas verificação incerta - URL: %s\n", voteNumber, resp.Request.URL.String())
		// Salva o body para debug apenas no primeiro voto com problema
		if voteNumber == 1 {
			fmt.Printf("DEBUG: Primeiros 500 caracteres da resposta:\n%s\n", bodyStr[:min(500, len(bodyStr))])
		}
		return nil
	}

	return fmt.Errorf("resposta inesperada: %d - URL: %s", resp.StatusCode, resp.Request.URL.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("  Auto Voter - Votação")
	fmt.Println("==========================================")
	fmt.Println()

	// Lê o link do formulário do arquivo config.txt
	configFile := "config.txt"
	formURL, err := readFormURLFromFile(configFile)
	if err != nil {
		fmt.Printf("❌ ERRO: %v\n", err)
		fmt.Println("\nCrie um arquivo 'config.txt' na mesma pasta com o link do formulário.")
		fmt.Println("Exemplo de conteúdo do config.txt:")
		fmt.Println("https://docs.google.com/forms/d/e/SEU_FORMULARIO_AQUI/viewform")
		return
	}

	fmt.Printf("📋 Link de votação carregado: %s\n", formURL)
	fmt.Println("  Votando  ")
	fmt.Println("  (Marcando: esta opção!)")
	fmt.Println("==========================================")
	fmt.Println()

	// Número de votos a serem enviados
	totalVotes := 1000000
	successCount := 0
	errorCount := 0

	fmt.Printf("Enviando %d votos...\n\n", totalVotes)

	for i := 1; i <= totalVotes; i++ {
		err := submitVote(i, formURL)
		if err != nil {
			fmt.Printf("✗ Erro no voto #%d: %v\n", i, err)
			errorCount++
		} else {
			successCount++
		}

		// Pausa entre os votos para não sobrecarregar o servidor
		// Ajuste o tempo conforme necessário
		if i < totalVotes {
			time.Sleep(2000 * time.Millisecond)
		}

		// Mostra progresso a cada 10 votos
		if i%10 == 0 {
			fmt.Printf("\n--- Progresso: %d/%d votos enviados ---\n\n", i, totalVotes)
		}
	}

	// Resumo final
	fmt.Println("\n==========================================")
	fmt.Println("  Resumo Final")
	fmt.Println("==========================================")
	fmt.Printf("Total de votos tentados: %d\n", totalVotes)
	fmt.Printf("✓ Votos bem-sucedidos:   %d\n", successCount)
	fmt.Printf("✗ Votos com erro:        %d\n", errorCount)
	fmt.Printf("Taxa de sucesso:         %.1f%%\n", float64(successCount)/float64(totalVotes)*100)
	fmt.Println("==========================================")
}
