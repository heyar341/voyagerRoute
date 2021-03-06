package simulsearch

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var simulSearchTpl *template.Template

func init() {
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/simul_search.html")).ParseGlob("templates/includes/*.html"))
}

func SimulSearchTpl(w http.ResponseWriter, req *http.Request) {
	msg := url.QueryEscape("エラーが発生しました。しばらく経ってからもう一度ご利用ください。")
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		log.Printf("Error whle gettibg data from context")
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	simulSearchTpl.ExecuteTemplate(w, "simul_search.html", data)
}
