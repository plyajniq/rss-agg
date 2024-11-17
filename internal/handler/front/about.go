package front

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func LoadAboutTemplate() *template.Template {
	base_tmpt := filepath.Join("..", "templates", "base.html")
	about_tmpt := filepath.Join("..", "templates", "content", "about.html")
	head_tmpt := filepath.Join("..", "templates", "core", "head.html")
	header_tmpt := filepath.Join("..", "templates", "core", "header.html")
	footer_tmpt := filepath.Join("..", "templates", "core", "footer.html")

	tmpl, err := template.ParseFiles(base_tmpt, about_tmpt, head_tmpt, header_tmpt, footer_tmpt)
	if err != nil {
		log.Fatal("Ошибка при загрузке шаблонов:", err)
	}
	return tmpl
}

func GetAbout(w http.ResponseWriter, r *http.Request) {
	tmpl := LoadAboutTemplate()
	data := map[string]interface{}{
		"About":  "Он работает!",
	}
	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Println("Ошибка при рендеринге шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
