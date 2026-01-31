package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bizflow/internal/ai"
)

// enableCORS allows your Svelte frontend to talk to this backend
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleRunAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input ai.BusinessInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Initialize and run
	notion := ai.NewNotionClient()
	output := ai.RunAgent(input, notion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/run-agent", handleRunAgent)

	// Wrap with CORS
	handler := enableCORS(mux)

	// 0.0.0.0 is the secret sauce for mobile/network access
	port := "8080"
	addr := "0.0.0.0:" + port

	fmt.Println("🚀 BizFlow Backend is LIVE")
	fmt.Printf("💻 Local Access:   http://localhost:%s\n", port)
	fmt.Printf("📱 Network Access: http://192.168.89.196:%s\n", port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
