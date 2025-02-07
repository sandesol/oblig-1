package main

import (
	"ASS/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Standard http server with reference to stubbed handler
	// If you want to adapt this, ensure to adjust path for compatibility with project
	http.HandleFunc("/", handlers.DiagHandler)

	// Naturally, you can introduce multiple handlers to emulate different data sources
	//http.HandleFunc("/species/occurrences", internal.StubHandlerOccurrences)

	log.Println("Running on port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
