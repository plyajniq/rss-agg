package api

import (
	"net/http"

	"rss-agg/internal/utils"
)

func Error(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusBadRequest, "Request goes wrong")
}
