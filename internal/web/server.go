package web

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wyll-io/dicomizer/internal/database"
	"github.com/wyll-io/dicomizer/internal/web/authentication"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
	webFuncs "github.com/wyll-io/dicomizer/internal/web/functions"
	"github.com/wyll-io/dicomizer/internal/web/home"
	"github.com/wyll-io/dicomizer/internal/web/patient"
)

var (
	internalCtx webContext.InternalValues
	templates   = webContext.TemplatesValues{}
)

func init() {
	templates["home"] = template.Must(
		template.New("home").
			Funcs(template.FuncMap{
				"map":   webFuncs.CreateMap,
				"isset": webFuncs.Isset,
			}).
			ParseFiles(
				"templates/common/base.html",
				"templates/common/header.html",
				"templates/common/footer.html",
				"templates/home/main.html",
				"templates/home/partials/add_patient.html",
				"templates/home/patient_list.html",
				"templates/home/partials/patient_row.html",
				"templates/home/partials/patient_table.html",
				"templates/home/partials/patient_row_edit.html",
				"templates/home/partials/text_input.html",
			),
	)
	templates["error"] = template.Must(
		template.New("error").
			ParseFiles(
				"templates/error.html",
				"templates/common/base.html",
				"templates/common/header.html",
				"templates/common/footer.html",
			),
	)
	templates["login"] = template.Must(
		template.New("login").
			ParseFiles(
				"templates/login.html",
				"templates/common/base.html",
				"templates/common/header.html",
				"templates/common/footer.html",
			),
	)
}

func RegisterHandlers(awsCfg aws.Config, dynamoDBTable string) http.Handler {
	internalCtx = webContext.InternalValues{
		DB: database.New(awsCfg, dynamoDBTable),
	}

	r := mux.NewRouter()

	home.Register(r)
	patient.Register(r)

	r.HandleFunc("/login", authentication.Handle).Methods("GET", "POST")
	r.HandleFunc("/logout", authentication.HandleLogout).Methods("GET")

	r.HandleFunc("/error", webError.Handle).Methods("GET")

	r.PathPrefix("/public").
		Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public")))).
		Methods("GET")

	return dispatchHandlers(r)
}

// dispatchHandlers sets up all the required http handlers
func dispatchHandlers(h http.Handler) http.Handler {
	h = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			internalCtx := context.WithValue(r.Context(), webContext.Internal, internalCtx)
			next.ServeHTTP(w, r.WithContext(context.WithValue(internalCtx, webContext.Templates, templates)))
		})
	}(h)
	h = handlers.LoggingHandler(os.Stdout, h)

	return authentication.AccessGuard(h)
}
