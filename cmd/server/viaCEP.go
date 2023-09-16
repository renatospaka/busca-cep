package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type viaCEP struct {
	CEP         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	UF          string
	IBGE        string
	GIA         string
	DDD         string
	SIAFI       string
}

func buscaViaCEP(cep string, ch chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*350)
	defer cancel()

	url := "http://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- err.Error()
	}
	// defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- err.Error()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- err.Error()
	}

	var c viaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		ch <- err.Error()
	}
	endereco := c.Logradouro + ", " + c.Bairro + ", " + c.Localidade + " - " + c.UF + ", " + c.CEP
	log.Printf("Via: %s\n", endereco)
	ch <- endereco
}
