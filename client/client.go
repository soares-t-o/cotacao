package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Usdbrl struct {
	Bid string `json:"bid"`
}

func main() {

	ctx, cancel := context.WithTimeout(
		context.Background(), 300*time.Millisecond,
	)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	var usd Usdbrl
	err = json.Unmarshal(body, &usd)

	if err != nil {
		panic(err)
	}

	saveQuote(usd)

	json.NewEncoder(os.Stdout).Encode(usd)
	fmt.Printf("Cotação do dólar: %s", usd.Bid)
}

func saveQuote(usd Usdbrl) {
	f, err := os.Create("./client/cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString(fmt.Sprintf("Dólar: {%s}\n", usd.Bid))

}
