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

// load `/posts` template
func LoadPostsTemplate() *template.Template {
	base_tmpt := filepath.Join("..", "templates", "base.html")
	posts_tmpt := filepath.Join("..", "templates", "content", "posts.html")
	head_tmpt := filepath.Join("..", "templates", "core", "head.html")
	header_tmpt := filepath.Join("..", "templates", "core", "header.html")
	footer_tmpt := filepath.Join("..", "templates", "core", "footer.html")

	tmpl, err := template.ParseFiles(base_tmpt, posts_tmpt, head_tmpt, header_tmpt, footer_tmpt)
	if err != nil {
		log.Fatal("Fail to load templates:", err)
	}
	return tmpl
}

// get posts by feed
func GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)

	feedIDStr := chi.URLParam(r, "feedID")
	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to parse feedID from URL: %v", err))
		return
	}
	feed, err := db.GetFeedByID(r.Context(), feedID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get feed by id: %v", err))
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

	data := map[string]interface{}{
		"Head":   feed.Name,
		"Posts":  posts,
		"Feed":   feed,
	}

	tmpl := LoadPostsTemplate()
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Println("Fail to render template:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
