package main

import (
	"fmt"
	"net/http"

	"github.com/plyajniq/rss-agg/internal/auth"
	"github.com/plyajniq/rss-agg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// run handler for authenticated user
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Fail to get API Key: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
