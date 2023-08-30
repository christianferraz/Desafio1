package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dolar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type USDBRL struct {
	ID         int `gorm:"primaryKey;autoIncrement:true"`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleCotacao)
	http.ListenAndServe(":8080", mux)
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo de log:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("error %v", err.Error())
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error %v", err.Error())
		panic(err)
	}
	defer resp.Body.Close()
	ContextTime(ctx, "Request site economia")
	if resp.StatusCode != http.StatusOK {
		log.Printf("error %v", err.Error())
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error %v", err.Error())
		return
	}
	dolar := Dolar{}

	err = json.Unmarshal(body, &dolar)
	if err != nil {
		log.Printf("error %v", err.Error())
		panic(err)
	}
	err = SaveDB(&dolar)
	if err != nil {
		log.Printf("error %v", err.Error())
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"Bid": dolar.USDBRL.Bid})
}

func SaveDB(dolar *Dolar) error {
	db, err := gorm.Open(sqlite.Open("go.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&USDBRL{})
	usdbrl := USDBRL{
		Code:       dolar.USDBRL.Code,
		Codein:     dolar.USDBRL.Codein,
		Name:       dolar.USDBRL.Name,
		High:       dolar.USDBRL.High,
		Low:        dolar.USDBRL.Low,
		VarBid:     dolar.USDBRL.VarBid,
		PctChange:  dolar.USDBRL.PctChange,
		Bid:        dolar.USDBRL.Bid,
		Ask:        dolar.USDBRL.Ask,
		Timestamp:  dolar.USDBRL.Timestamp,
		CreateDate: dolar.USDBRL.CreateDate,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db.WithContext(ctx).Create(&usdbrl)
	ContextTime(ctx, "DB")
	return nil
}

func ContextTime(ctx context.Context, t string) {
	select {
	case <-ctx.Done():
		log.Printf("Time exceeded in " + t)
	case <-time.After(time.Millisecond * 15):
		println("success")
	}
}
