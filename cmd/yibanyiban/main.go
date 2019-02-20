package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kyeett/yibanyiban"
)

func main() {
	portFlag := flag.String("port", "8080", "Port number for server")
	flag.Parse()

	srv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + *portFlag,
	}
	http.HandleFunc("/validate", yibanyiban.ValidateIBANHandler)

	fmt.Printf("Serving IBAN validation service on :%s\n", *portFlag)
	log.Fatal(srv.ListenAndServe())
}
