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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)
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

	var c apiCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		ch <- err.Error()
	}
	endereco := c.Address + ", " + c.District + ", " + c.City + " - " + c.State + ", " + c.Code
	ch <- endereco
}