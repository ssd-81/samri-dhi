package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"fd-credit-score/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handlers.GetHealth)
	mux.HandleFunc("POST /api/score", handlers.CalculateScore)
	mux.HandleFunc("GET /api/personas", handlers.ListPersonas)
	mux.HandleFunc("GET /api/personas/{id}/score", handlers.GetPersonaScore)

	wd, _ := os.Getwd()
	fs := http.Dir(wd + "/web")
	mux.Handle("/", http.FileServer(fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("FD Credit Score running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}