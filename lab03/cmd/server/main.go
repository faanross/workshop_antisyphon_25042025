package main

import (
	"fmt"
	"github.com/faanross/orlokC2/internal/router"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const serverAddr = "127.0.0.1"
const serverPort = "7777"

func main() {

	r := chi.NewRouter()

	router.SetupRoutes(r)

	serverAddrPort := fmt.Sprintf("%s:%s", serverAddr, serverPort)

	log.Printf("Server listening on %s\n", serverAddrPort)

	go func() {
		err := http.ListenAndServe(serverAddrPort, r)
		if err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()

	select {}

}
