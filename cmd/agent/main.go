package main

import (
    "net/http"
    "fmt"
    "encoding/json"

    "bizflow/internal/ai"
)

func main() {

    http.HandleFunc("/run-agent", func(w http.ResponseWriter, r *http.Request) {

        if r.Method != http.MethodPost {
            http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
            return
        }

        var input ai.BusinessInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        output := ai.RunAgent(input)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(output)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK")
    })

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}