package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AddressViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type AddressBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {

	cep := "01001000"
	getFastServiceCep(cep)

}

func getFastServiceCep(cep string) {
	ch := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	go BuscaCepBrasilApi(ctx, ch, cep)
	go BuscaCepViaCep(ctx, ch, cep)

	select {
	case resp := <-ch:
		println(resp)
	case <-ctx.Done():
		println("timeout")
	}
}

func BuscaCepViaCep(ctx context.Context, ch chan string, cep string) (*AddressViaCep, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data AddressViaCep
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	ch <- fmt.Sprintf("Service: ViaCep | CEP: %s - Cidade: %s - UF: %s - Street: %s - Bairro: %s", data.Cep, data.Localidade, data.Uf, data.Localidade, data.Bairro)
	return &data, nil
}

func BuscaCepBrasilApi(ctx context.Context, ch chan string, cep string) (*AddressBrasilApi, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data AddressBrasilApi
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	ch <- fmt.Sprintf("Service: BrasilAPI | CEP: %s - Cidade: %s - UF: %s - Street: %s - Neighborhood: %s", data.Cep, data.City, data.State, data.Street, data.Neighborhood)
	return &data, nil
}
