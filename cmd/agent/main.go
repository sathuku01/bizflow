package main

import (
	"biz-flow/web"
	"log"
	"net/http"
)

func main() {
	// Set up the API endpoint
	http.HandleFunc("/api/consult", web.ConsultationHandler)

	// Serve the static frontend
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	port := ":8080"
	log.Printf("Starting BizFlow server on http://localhost%s\n", port)

	// Start the server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
