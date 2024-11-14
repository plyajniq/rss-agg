package handler

import (
	"net/http"

	"rss-agg/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
