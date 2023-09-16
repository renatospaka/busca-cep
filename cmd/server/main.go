package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/{cep}", buscaCEP)

	log.Println("listening at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func buscaCEP(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	log.Printf("CEP: %s\n", cep)
	if cep == "" || len(cep) != 9 {	
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "cep not provided", http.StatusUnprocessableEntity)
	}

	api := make(chan string)
	via := make(chan string)

	// api CEP
	go buscaAPICEP(cep, api)

	// via CEO
	go buscaViaCEP(cep, via)

	// recebe o 1º resultado que chega
	w.Header().Set("Content-Type", "application/json")
	for {
		select {
		case endereco1 := <- api:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("API API:" + endereco1))
			// log.Printf("API API: %s\n", endereco1)
			
		case endereco2 := <- via:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Via API: " + endereco2))
			// log.Printf("Via API: %s\n", endereco2)

		case <-time.After(time.Second * 1):
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("timeout after 1 second"))
			// log.Printf("timeout after 1 second")
		}		
	}
}
