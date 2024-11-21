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

type NewUser struct {
	Name string `json:"name"`
}

// @Summary Create user
// @Description to create a new user's APIKey
// @Tags users
// @Accept json
// @Produce json
// @Param name body NewUser true "Name"
// @Success 201 {object} utils.User
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)

	decoder := json.NewDecoder(r.Body)

	params := NewUser{}

	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to parse JSON: %v", err))
		return
	}

	newUser, err := db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to create user: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, utils.DatabaseUserToUser(newUser))
}

// hander to get existing user by apy key
func GetMyUserData(w http.ResponseWriter, r *http.Request) {
	user := ctx.GetUserContext(r)
	utils.RespondWithJSON(w, http.StatusOK, utils.DatabaseUserToUser(user))
}

// hander to get existing user by apy key
func GetPostsForUser(w http.ResponseWriter, r *http.Request) {
	user := ctx.GetUserContext(r)
	db := ctx.GetDBContext(r)

	posts, err := db.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get posts for user: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.DatabasePostsToPosts(posts))
}
