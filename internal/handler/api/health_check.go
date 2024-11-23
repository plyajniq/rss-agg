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

// @Summary Health Check
// @Description to check server status
// @Tags service support
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /healthz [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(
		w,
		http.StatusOK,
		HealthResponse{
			Status:      "Server is running",
			CurrentTime: time.Now().String(),
		})
}
