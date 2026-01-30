package web

import (
	"biz-flow/internal/ai"
	"biz-flow/internal/core"
	"biz-flow/internal/handler"
	"encoding/json"
	"log"
	"net/http"
)

// ConsultationHandler handles requests for a new marketing consultation.
func ConsultationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var businessInput core.BusinessInput
	if err := json.NewDecoder(r.Body).Decode(&businessInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In a real application, we would call businessInput.Validate() here
	// once it's implemented by Developer A.

	// For now, we use the MockLLMClient. When Dev B is done,
	// this would be replaced with the real client, e.g., ai.NewOpenAIClient().
	client := ai.NewMockLLMClient()

	result, err := handler.GenerateConsultation(businessInput, client)
	if err != nil {
		log.Printf("Error generating consultation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}