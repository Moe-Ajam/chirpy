package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting...")
	mux := http.NewServeMux()
	rh := http.RedirectHandler("http://localhost", 307)

	mux.Handle("/", rh)

	log.Println("Listening...")

	http.ListenAndServe(":8080", mux)
}
