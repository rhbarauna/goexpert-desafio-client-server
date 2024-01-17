package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Iniciando o servidor")
	http.HandleFunc("/cotacao", filtraCotacaoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func filtraCotacaoHandler(res http.ResponseWriter, req *http.Request) {
	select {
	case <-req.Context().Done():
		log.Println("Request cancelada pelo cliente")
		return
	default:
	}

	cotacao, err := buscaCotacao(res)
	if err != nil {
		return
	}

	err = persistIntoDB(*cotacao)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Write([]byte(cotacao.Usdbrl.Bid))
}

func buscaCotacao(res http.ResponseWriter) (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	api_response, err := http.DefaultClient.Do(req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	defer api_response.Body.Close()

	body, err := io.ReadAll(api_response.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return &cotacao, nil
}

func prepareDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS cocatao (
		id TEXT NOT NULL,
		code TEXT NOT NULL,
		codein TEXT NOT NULL,
		name TEXT NOT NULL,
		high TEXT NOT NULL,
		low TEXT NOT NULL,
		varBid TEXT NOT NULL,
		pctChange TEXT NOT NULL,
		bid TEXT NOT NULL,
		ask TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		create_date TEXT NOT NULL
	);`)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func persistIntoDB(c Cotacao) error {
	db, err := sql.Open("sqlite3", "db_cotacao.sqlite3")

	if err != nil {
		fmt.Println("Erro ao preparar o banco")
		return err
	}
	defer db.Close()

	err = prepareDB(db)
	if err != nil {
		fmt.Println("Erro ao preparar o banco")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("INSERT INTO cocatao (id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		fmt.Printf("Erro ao preparar a sql %v \n", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		uuid.New().String(),
		c.Usdbrl.Code,
		c.Usdbrl.Codein,
		c.Usdbrl.Name,
		c.Usdbrl.High,
		c.Usdbrl.Low,
		c.Usdbrl.VarBid,
		c.Usdbrl.PctChange,
		c.Usdbrl.Bid,
		c.Usdbrl.Ask,
		c.Usdbrl.Timestamp,
		c.Usdbrl.CreateDate)

	if err != nil {
		fmt.Println("Erro ao persistir dados no database")
		return err
	}

	return nil
}

type USDBRL struct {
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
}

type Cotacao struct {
	Usdbrl USDBRL `json:"USDBRL"`
}
