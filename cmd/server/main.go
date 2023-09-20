package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewMux()
	mux.Route("/busca", func(r chi.Router) {
		r.Get("/{cep}", buscaCEP)
	})

	log.Println("listening at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func buscaCEP(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()
	cep := chi.URLParam(r, "cep")
	// if cep == "" || (len(cep) < 8 || len(cep) > 9) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	http.Error(w, "cep not provided", http.StatusUnprocessableEntity)
	// }

	b := []byte(cep)
	rn := regexp.MustCompile(`^\d{8}$`)
	if rn.Match(b) {
		cep = cep[0:5] + "-" + cep[5:8]
		b = []byte(cep)
	} 
	
	rt := regexp.MustCompile(`^\d{5}[-]\d{3}$`)
	if !rt.Match(b) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "cep invalid", http.StatusBadRequest)
		return
	}
	log.Printf("searching for %s\n", cep)

	api := make(chan string)
	via := make(chan string)

	// api CEP
	go buscaAPICEP(cep, api)

	// // via CEP
	go buscaViaCEP(cep, via)

	// formata o header
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	
	// recebe o 1º resultado que chega
	select {
	case endereco1 := <- api:
		log.Printf("Duration: %s", time.Since(startAt).String())
		response := "API API:" + endereco1
		encoder.Encode(&response)

	case endereco2 := <- via:
		log.Printf("Duration: %s", time.Since(startAt).String())
		response := "Via API:" + endereco2
		encoder.Encode(&response)

	case <-time.After(time.Second * 1):
		log.Printf("Duration: %s", time.Since(startAt).String())
		w.WriteHeader(http.StatusRequestTimeout)
		encoder.Encode("timeout after 1 second")
	}
}
