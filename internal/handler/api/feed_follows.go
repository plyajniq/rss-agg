package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ctx "rss-agg/internal/context"
	"rss-agg/internal/database"
	"rss-agg/internal/utils"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// handler to create a new feed
func CreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to parse JSON: %v", err))
		return
	}

	feedFollow, err := db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to create feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, utils.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func GetFeedFollows(w http.ResponseWriter, r *http.Request) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	feedFollows, err := db.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get feed follows by user: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func DeleteFeedFollows(w http.ResponseWriter, r *http.Request) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail parse feed follow id: %v", err))
		return
	}

	err = db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail delete feed follow: %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, struct{}{})
}
