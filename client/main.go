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
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo de log:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(nil)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(nil)
	}
	defer res.Body.Close()
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()
	select {
	case <-ctx.Done():
		log.Printf("Time exceeded")
	case <-time.After(time.Millisecond * 200):
		_, err = io.Copy(file, res.Body)
		if err != nil {
			fmt.Println("Erro ao copiar o corpo da resposta:", err)
		}
	}
}
