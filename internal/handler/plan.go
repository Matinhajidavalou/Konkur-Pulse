package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"konkur-pulse/internal/repository"
	"konkur-pulse/internal/service"
)

// GetPlanHandler returns an HTTP handler for generating daily study plans.
// It fetches upstream data via MazeClient and allocates time accordingly.
func GetPlanHandler(mazeClient *repository.MazeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID := r.PathValue("id")
		dateStr := r.URL.Query().Get("date")

		loc, _ := time.LoadLocation("Asia/Tehran")
		_, err := time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
			return
		}

		// Fetch external results (utilizes internal cache if available)
		subjects, err := mazeClient.FetchResults(studentID, dateStr)
		if err != nil {
			http.Error(w, `{"error": "Upstream error: `+err.Error()+`"}`, http.StatusBadGateway)
			return
		}

		if len(subjects) == 0 {
			http.Error(w, `{"error": "No valid subjects"}`, http.StatusNotFound)
			return
		}

		// Transform the fetched struct slice into a map for the allocation service
		subjectScores := make(map[string]int)
		for _, s := range subjects {
			subjectScores[s.Subject] = s.Traz
		}

		allocatedPlan, err := service.AllocateStudyTime(subjectScores, 240, 8000, 30)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"student_id": studentID,
			"plan":       allocatedPlan,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
