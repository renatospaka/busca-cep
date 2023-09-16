package main

import (
	// "log"
	"fmt"
	"log"
	// "net/http"
)

func main() {
	// http.HandleFunc("/", buscaCEP)

	cep := "04165-030"

	e1, err := buscaAPICEP(cep)
	if err != nil {
		log.Fatalf("Error on APICEP: %v\n", err)
	}
	fmt.Printf("APICEP: %s \n", e1)

	e2, err := buscaViaCEP(cep)
	if err != nil {
		log.Fatalf("Error on ViaCEP: %v\n", err)
	}
	fmt.Printf("viaCEP: %s \n", e2)

	// log.Println("listening at http://localhost:3000")
	// http.ListenAndServe(":3000", nil)
}

// func buscaCEP(w http.ResponseWriter, r *http.Request) {

// }
