package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("CEP Obrigatorio para consulta")
		panic("CEP Obrigatorio para consulta")
	}

	cep := os.Args[1]

	uriViacep := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	uriBrazilapi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	channelViacep := make(chan string)
	channelBrazilapi := make(chan string)

	go execute(uriViacep, channelViacep)
	go execute(uriBrazilapi, channelBrazilapi)

	select {
	case msg := <-channelViacep:
		fmt.Print(msg)

	case msg := <-channelBrazilapi:
		fmt.Print(msg)

	case <-time.After(time.Second):
		fmt.Printf("Nenhuma das apis respondeu no tempo configurado")
	}

}

func execute(uri string, ch chan string) {
	req, err := http.Get(uri)
	if err != nil {
		fmt.Printf("Erro ao realizar consulta na url %s", uri)
		panic(err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Erro ao realizar consulta na url %s", uri)

		panic(err)
	}
	ch <- string(res)
}
