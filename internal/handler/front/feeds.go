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

func LoadFeedsTempalate() *template.Template {
	base_tmpt := filepath.Join("..", "templates", "base.html")
	feeds_tmpt := filepath.Join("..", "templates", "content", "feeds.html")
	head_tmpt := filepath.Join("..", "templates", "core", "head.html")
	header_tmpt := filepath.Join("..", "templates", "core", "header.html")
	footer_tmpt := filepath.Join("..", "templates", "core", "footer.html")

	tmpl, err := template.ParseFiles(base_tmpt, feeds_tmpt, head_tmpt, header_tmpt, footer_tmpt)
	if err != nil {
		log.Fatal("Ошибка при загрузке шаблонов:", err)
	}
	return tmpl
}

func GetTopFeeds(w http.ResponseWriter, r *http.Request) {
	db := ctx.GetDBContext(r)

	feeds, err := db.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fail to get all feeds: %v", err))
		return
	}
	tmpl := LoadFeedsTempalate()

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
