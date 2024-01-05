package home

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["home"]
	if err := tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"LoggedIn": r.Context().Value(webContext.User).(bool),
		"Title":    "Dicomizer",
		"Patients": []dao.Patient{
			{
				UUID:      "some_uuid",
				Firstname: "Antoine",
				Lastname:  "Langlois",
				Filters:   "some,filters",
				Studies:   []dao.DCMImage{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Time{},
				DeletedAt: time.Time{},
			},
			{
				UUID:      "another_uuid",
				Firstname: "John",
				Lastname:  "Doe",
				Filters:   "other,filters",
				Studies:   []dao.DCMImage{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Time{},
				DeletedAt: time.Time{},
			},
			{
				UUID:      "third_uuid",
				Firstname: "Jane",
				Lastname:  "Smith",
				Filters:   "third,filters",
				Studies:   []dao.DCMImage{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Time{},
				DeletedAt: time.Time{},
			},
			{
				UUID:      "fourth_uuid",
				Firstname: "Michael",
				Lastname:  "Johnson",
				Filters:   "fourth,filters",
				Studies:   []dao.DCMImage{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Time{},
				DeletedAt: time.Time{},
			},
		},
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
