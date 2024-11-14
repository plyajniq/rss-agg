package ctx

import (
	"net/http"

	"rss-agg/internal/database"
)

func GetUserContext(r *http.Request) database.User {
	return r.Context().Value("user").(database.User)
}

func GetDBContext(r *http.Request) *database.Queries {
	return r.Context().Value("db").(*database.Queries)
}
