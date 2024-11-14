package middleware

import (
	"context"
	"fmt"
	"net/http"

	"rss-agg/internal/auth"
	"rss-agg/internal/database"
	"rss-agg/internal/utils"
)

// check authentication for request
func BasicAuth(db *database.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, err := auth.GetAPIKey(r.Header)
			if err != nil {
				utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Fail to get API Key: %v", err))
				return
			}

			user, err := db.GetUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get user: %v", err))
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


