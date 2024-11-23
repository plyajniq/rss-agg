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

type NewFeedFollow struct {
	FeedID uuid.UUID `json:"feed_id"`
}

// @Summary Create Feed Follow
// @Description to create a new feed follow
// @Tags feed follows
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "APIKey" example(ApiKey $token)
// @Param feed_id body NewFeedFollow true "FeedID"
// @Success 201 {object} utils.FeedFollow
// @Failure 400 {object} utils.ErrResponse
// @Router /feed_follows [post]
func CreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	decoder := json.NewDecoder(r.Body)
	params := NewFeedFollow{}

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

// @Summary Get Feed Follows
// @Description to get feed follows
// @Tags feed follows
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "APIKey" example(ApiKey $token)
// @Success 200 {array} utils.FeedFollow
// @Failure 400 {object} utils.ErrResponse
// @Router /feed_follows [get]
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

// @Summary Delete Feed Follow
// @Description to delete feed follow (unfollow)
// @Tags feed follows
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "APIKey" example(ApiKey $token)
// @Param feed_follow_id path string true "Feed Follow ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Router /feed_follows/{feed_follow_id} [delete]
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
