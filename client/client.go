package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ResponseCotacao struct {
	Bid string `json:"bid"`
}

func main(){
	ctx,cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	bid, err := GetCotacao(ctx)
	if err != nil {
		log.Println("Erro ao obter cotação:", err)
	}
	if err := SaveFile(bid); err != nil {
		log.Println("Erro ao salvar cotação no arquivo:", err)
	}
}
func GetCotacao(ctx context.Context) (string, error){
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var responsecotacao ResponseCotacao
	if err := json.NewDecoder(res.Body).Decode(&responsecotacao); err != nil {
		return "", err
	}
	return responsecotacao.Bid, nil
}

func SaveFile(bid string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dolar: %s\n", bid))
	
	return err
}