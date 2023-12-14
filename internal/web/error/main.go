package error

import (
	"fmt"
	"net/http"
	"strconv"

	webContext "github.com/wyll-io/dicomizer/internal/web/context"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	qc := r.URL.Query().Get("code")
	c, err := strconv.Atoi(qc)
	if err != nil {
		fmt.Printf("error while parsing error code \"%s\": %v", qc, err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["error"]
	if err := tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Code":     c,
		"Msg":      r.URL.Query().Get("message"),
		"LoggedIn": r.Context().Value(webContext.User).(bool),
		"Title":    "Dicomizer - Error",
	}); err != nil {
		panic(fmt.Sprintf("error while loading error template: %v", err))
	}
}
