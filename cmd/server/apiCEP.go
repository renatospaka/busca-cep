package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type apiCEP struct {
	Code       string
	State      string
	City       string
	District   string
	Address    string
	Status     int
	OK         bool
	StatusText string
}

func buscaAPICEP(cep string, ch chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*750)
	defer cancel()

	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- err.Error()
	}
	// defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- err.Error()
	}
	if resp == nil {
		ch <- "sem resposta"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- err.Error()
	}
	if body == nil {
		ch <- "sem body"
	}

	var c apiCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		ch <- err.Error()
	}
	endereco := c.Address + ", " + c.District + ", " + c.City + " - " + c.State + ", " + c.Code
	// log.Printf("API: %s\n", endereco)
	ch <- endereco
}
