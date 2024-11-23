package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ctx "rss-agg/internal/context"
	"rss-agg/internal/database"
	"rss-agg/internal/utils"

	"github.com/google/uuid"
)

type NewFeed struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// @Summary Create Feed
// @Description to create a new feed
// @Tags feeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "APIKey" example(ApiKey $token)
// @Param name body NewFeed true "Name"
// @Param url body NewFeed true "URL"
// @Success 201 {object} utils.Feed
// @Failure 400 {object} utils.ErrResponse
// @Router /feeds [post]
func CreateFeed(w http.ResponseWriter, r *http.Request) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	decoder := json.NewDecoder(r.Body)
	params := NewFeed{}

	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to parse JSON: %v", err))
		return
	}

	feed, err := db.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to create feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, utils.DatabaseFeedToFeed(feed))
}

// @Summary Get Feeds
// @Description to get feeds
// @Tags feeds
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "APIKey" example(ApiKey $token)
// @Success 200 {array} utils.Post
// @Failure 400 {object} utils.ErrResponse
// @Router /feeds [get]
func GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)
	feeds, err := db.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get all feeds: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.DatabaseFeedsToFeeds(feeds))

}
