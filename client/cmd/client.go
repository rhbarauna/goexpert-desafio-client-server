package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err == context.DeadlineExceeded {
			fmt.Println("Tempo de execução excedido")
		}
		return
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Falha ao construir a requisição: %s\n", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Falha ao obter a cotação atual: %s\n", err.Error())
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Falha ao ler o corpo da resposta: %s\n", err.Error())
		return
	}

	err = createCotacaoFile(body)

	if err != nil {
		return
	}

	log.Printf("Cotação salva com sucesso!")
}

func createCotacaoFile(content []byte) error {
	_, err := os.Stat("./cotacao.txt")

	if os.IsNotExist(err) {
		_, err = os.Create("./cotacao.txt")
	}

	if err != nil {
		log.Printf("Falha ao verificar o arquivo: %s\n", err.Error())
		return err
	}

	err = os.WriteFile("./cotacao.txt",
		[]byte(fmt.Sprintf("Dólar: %s", content)),
		0644,
	)

	if err != nil {
		log.Printf("Falha ao escrever a cotação no arquivo: %s\n", err.Error())
		return err
	}

	return nil
}
