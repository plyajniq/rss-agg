package api

import (
	"net/http"
	"time"

	"rss-agg/internal/utils"
)

type HealthResponse struct {
	Status      string
	CurrentTime string
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(
		w,
		http.StatusOK,
		HealthResponse{
			Status:      "Server is running",
			CurrentTime: time.Now().String(),
		})
}
