package handlers

import (
	"encoding/json"
	"fd-credit-score/internal/personas"
	"fd-credit-score/internal/scoring"
	"net/http"
)

func CalculateScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var history scoring.DepositHistory
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := scoring.CalculateScore(history)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func ListPersonas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	all, err := personas.GetAll()
	if err != nil {
		http.Error(w, "Failed to load personas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func GetPersonaScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Persona ID required", http.StatusBadRequest)
		return
	}

	persona, err := personas.GetByID(id)
	if err != nil {
		http.Error(w, "Persona not found", http.StatusNotFound)
		return
	}

	history := scoring.DepositHistory{
		UserID:  persona.ID,
		Name:    persona.Name,
		Age:     persona.Age,
		City:    persona.City,
		Deposits: persona.Deposits,
	}

	result := scoring.CalculateScore(history)
	summary := scoring.ComputePersonaSummary(history, persona.Name, persona.Age, persona.Occupation, persona.City)

	response := map[string]interface{}{
		"persona":      summary,
		"total_score":  result.TotalScore,
		"score_band":   result.ScoreBand,
		"projected_score": result.ProjectedScore,
		"cibil_equivalent": result.CIBILEquivalent,
		"components":   result.Components,
		"credit_products": result.CreditProducts,
		"improvements": result.Improvements,
		"insights":     result.Insights,
		"pattern_detected": result.PatternDetected,
		"peer_percentile": result.PeerPercentile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}