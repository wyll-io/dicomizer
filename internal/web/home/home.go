package home

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iCtxVal := ctx.Value(webContext.Internal).(webContext.InternalValues)
	patients, err := iCtxVal.DB.GetPatientsInfo(ctx)
	if err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	for i := range patients {
		// * expensive operation but sufficient for now
		new := patients[i]
		new.PK = strings.Replace(new.PK, "PATIENT#", "", 1)
		patients[i] = new
	}

	tmpl := ctx.Value(webContext.Templates).(webContext.TemplatesValues)["home"]
	if err := tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"LoggedIn": ctx.Value(webContext.User).(bool),
		"Title":    "Dicomizer",
		"Patients": patients,
	}); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func Register(r *mux.Router) {
	r.HandleFunc("/", handleHome).Methods("GET")
	r.HandleFunc("/partials/patient-row", handlePatientRow).Methods("GET")
	r.HandleFunc("/partials/patient-row-edit", handlePatientRowEdit).Methods("GET")
	r.HandleFunc("/validate/input/{id}", handleValidateSingleInput).Methods("POST")
}
