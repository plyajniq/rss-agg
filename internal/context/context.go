package ctx

import (
	"net/http"

	"rss-agg/internal/database"
)

// get user entity from context
func GetUserContext(r *http.Request) database.User {
	return r.Context().Value("user").(database.User)
}

// get database connection from context
func GetDBContext(r *http.Request) *database.Queries {
	return r.Context().Value("db").(*database.Queries)
}
