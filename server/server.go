package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
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

func main() {

	http.HandleFunc("/cotacao", GetQuote)
	http.ListenAndServe(":8080", nil)
}

func GetQuote(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(
		context.Background(), 200*time.Millisecond,
	)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var quote Quote

	err = json.Unmarshal(body, &quote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = Save(quote)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao salvar no banco de dados"))
		return
	}
	json.NewEncoder(w).Encode(quote.Usdbrl)

}

func Save(quote Quote) error {

	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Millisecond,
	)
	defer cancel()

	db, err := sql.Open("sqlite3", "./server/databsase.db")

	if err != nil {
		return err
	}

	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS quotes ( bid varchar(255) )")

	_, err = db.ExecContext(ctx,
		"INSERT INTO quotes ( bid ) VALUES ( $1 )",
		quote.Usdbrl.Bid,
	)

	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT * FROM quotes")

	if err != nil {
		return err
	}
	var quotes []Quote

	for rows.Next() {
		q := Quote{}
		err = rows.Scan(&q.Usdbrl.Bid)
		if err != nil {
			return err
		}
		quotes = append(quotes, q)
	}
	for _, q := range quotes {
		fmt.Printf("value bid: %s\n", q.Usdbrl.Bid)
	}
	return nil
}
