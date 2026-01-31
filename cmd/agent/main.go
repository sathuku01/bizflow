package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bizflow/internal/ai"
)

func main() {
	http.HandleFunc("/run-agent", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Recovered from panic:", rec)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		var input ai.BusinessInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		// Initialize Notion client
		notion := ai.NewNotionClient()
		output := ai.RunAgent(input, notion)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(output); err != nil {
			log.Println("failed to encode response:", err)
		}
	})

	fmt.Println("Agent server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
