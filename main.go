package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting...")
	mux := http.NewServeMux()
	server := http.Server{}
	handler := server.Handler
}
