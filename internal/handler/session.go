package handler

import (
	"encoding/json"
	"net/http"

	"konkur-pulse/internal/repository"
)

// SessionRequest defines the payload expected for saving a study session.
type SessionRequest struct {
	Date            string `json:"date"`
	Subject         string `json:"subject"`
	DurationMinutes int    `json:"duration_minutes"`
	ClientSessionID string `json:"client_session_id"` // Used as an idempotency key
}

// SaveSessionHandler returns an HTTP handler to persist study sessions.
// It ensures idempotent operations relying on ClientSessionID.
func SaveSessionHandler(repo *repository.SessionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID := r.PathValue("id")

		var req SessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		err := repo.SaveSession(studentID, req.Date, req.Subject, req.DurationMinutes, req.ClientSessionID)
		if err != nil {
			http.Error(w, `{"error": "Failed to save session"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Session saved successfully"}`))
	}
}
