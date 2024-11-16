package front

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	ctx "rss-agg/internal/context"
	"rss-agg/internal/utils"
	"text/template"
)

func GetTopFeeds(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)

	feeds, err := db.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get all feeds: %v", err))
		return
	}
	base_tmpt := filepath.Join("..", "templates", "base.html")
	feeds_tmpt := filepath.Join("..", "templates", "content", "feeds.html")

	tmpl, err := template.ParseFiles(base_tmpt, feeds_tmpt)
	if err != nil {
		log.Fatal("Ошибка при загрузке шаблонов:", err)
	}

	data := map[string]interface{}{
		"Header":  "Ахуеть, работает!",
		"Head":   "Что-то тут!",
		"Special": "Специально тут!",
		"Feeds":   feeds,
	}
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Println("Ошибка при рендеринге шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
