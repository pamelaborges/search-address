package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type VIACepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("CEP Obrigatorio para consulta")
		panic("CEP Obrigatorio para consulta")
	}

	cep := os.Args[1]

	viaCepResponse := VIACepResponse{}
	brazilAPIResponse := BrasilAPIResponse{}

	uriViacep := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	uriBrazilapi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	channelViacep := make(chan interface{})
	channelBrazilapi := make(chan interface{})

	go execute(uriViacep, channelViacep, &viaCepResponse)
	go execute(uriBrazilapi, channelBrazilapi, &brazilAPIResponse)

	select {
	case msg := <-channelViacep:
		fmt.Printf("A API mais rapida foi a VIACep, e o retorno foi %s", msg)

	case msg := <-channelBrazilapi:
		fmt.Printf("A API mais rapida foi a BrazilAPI, e o retorno foi %s", msg)

	case <-time.After(time.Second):
		fmt.Printf("Nenhuma das apis respondeu no tempo configurado")
	}

}

func execute(uri string, ch chan interface{}, responseInterface interface{}) {
	req, err := http.Get(uri)
	if err != nil {
		fmt.Printf("Erro ao realizar consulta na url %s", uri)
		panic(err)
	}
	defer req.Body.Close()

	err = json.NewDecoder(req.Body).Decode(&responseInterface)
	if err != nil {
		fmt.Printf("Erro ao decodificar JSON da URL %s", uri)
		panic(err)
	}

	ch <- responseInterface
}
