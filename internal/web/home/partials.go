package home

import (
	"net/http"

	"github.com/gorilla/mux"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

func handlePatientRowEdit(w http.ResponseWriter, r *http.Request) {
	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["home"]
	if err := tmpl.ExecuteTemplate(w, "partials/patient_row_edit", map[string]string{
		"Fullname": r.URL.Query().Get("fullname"),
		"Filters":  r.URL.Query().Get("filters"),
		"PK":     r.URL.Query().Get("pk"),
	}); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func handlePatientRow(w http.ResponseWriter, r *http.Request) {
	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["home"]
	if err := tmpl.ExecuteTemplate(w, "partials/patient_row", map[string]string{
		"Fullname": r.URL.Query().Get("fullname"),
		"Filters":  r.URL.Query().Get("filters"),
		"PK":     r.URL.Query().Get("pk"),
	}); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func handleValidateSingleInput(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	v := r.FormValue(id)

	d := map[string]interface{}{
		"ID":    id,
		"Value": v,
		"OK":    true,
	}

	switch id {
	case "fullname":
		d["Autocomplete"] = "name"
		d["Label"] = "Identité"
		d["Placeholder"] = "John Smith"
		if len(v) < 4 {
			d["Error"] = "L'identité doit contenir au moins 4 caractères"
			d["OK"] = false
		}
	case "filters":
		d["Autocomplete"] = "off"
		d["Label"] = "Filtres"
		d["Placeholder"] = "Filtres"
		if len(v) < 2 {
			d["Error"] = "Les filtres doivent contenir au moins 2 caractères"
			d["OK"] = false
		}
	}

	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["home"]
	if err := tmpl.ExecuteTemplate(w, "partials/text_input", d); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
	}
}
