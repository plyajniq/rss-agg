package front

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	ctx "rss-agg/internal/context"
	"rss-agg/internal/database"
	"rss-agg/internal/utils"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)

	feedIDStr := chi.URLParam(r, "feedID")
	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail parse feed id: %v", err))
		return
	}
	feed, err := db.GetFeedByID(r.Context(), feedID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail parse feed id: %v", err))
		return
	}

	posts, err := db.GetPostsForFeed(r.Context(), database.GetPostsForFeedParams{
		FeedID: feedID,
		Limit:  10,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get all posts: %v", err))
		return
	}
	base_tmpt := filepath.Join("..", "templates", "base.html")
	posts_tmpt := filepath.Join("..", "templates", "content", "posts.html")

	tmpl, err := template.ParseFiles(base_tmpt, posts_tmpt)
	if err != nil {
		log.Fatal("Ошибка при загрузке шаблонов:", err)
	}

	data := map[string]interface{}{
		"Header": "Ахуеть, работает!",
		"Head":   "Wow. HEAD!",
		"Posts":  posts,
		"Feed":   feed,
	}
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Println("Ошибка при рендеринге шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
