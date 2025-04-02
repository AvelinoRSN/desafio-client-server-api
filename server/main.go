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

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}
var db *sql.DB

func main() {
	initDB()
	http.HandleFunc("/cotacao", getCotacaoHandler)
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func initDB(){
	var err error
	db, err = sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`) 
	if err != nil {
		log.Fatal(err)
	}
}

func GetPrice(ctx context.Context) (string, error) {
	//time.Sleep(400 * time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil{
		return "", err
	}

	var cotacao Cotacao
	if err := json.Unmarshal(body, &cotacao); err != nil {
		return "", err
	}
	return cotacao.USDBRL.Bid, nil
}

func SaveBD(ctx context.Context, bid string) error {
	stmt, err := db.Prepare("INSERT INTO cotacao (bid) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, bid)
	return err
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
//	time.Sleep(300 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	bid, err := GetPrice(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		log.Println("Erro na API de cambio: ", err)
		return
	}
	
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()
	
	if err := SaveBD(dbCtx, bid); err != nil {
		http.Error(w, "Erro ao salvar cotação no banco de dados", http.StatusInternalServerError)
		log.Println("Erro ao salvar cotação no banco de dados: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": bid}) //retorna o bid
}