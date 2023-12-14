package patient

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	dao "github.com/wyll-io/dicomizer/internal/DAO"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

func delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: remove this when dynamodb is provisioned
	if false {
		iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
		if err := iCtxV.DB.DeletePatient(r.Context(), mux.Vars(r)["uuid"]); err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	search := r.FormValue("search")
	ps := []dao.Patient{
		{
			UUID:      "some_uuid",
			Firstname: "Antoine",
			Lastname:  "Langlois",
			Filters:   "some,filters",
			Studies:   []dao.Study{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
		{
			UUID:      "another_uuid",
			Firstname: "John",
			Lastname:  "Doe",
			Filters:   "other,filters",
			Studies:   []dao.Study{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}

	if false {
		iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
		_, err := iCtxV.DB.GetPatient(r.Context(), dao.SearchPatient{
			Firstname: search,
			Lastname:  search,
		}, false)
		if err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		d := map[string]interface{}{
			"Patients": ps,
		}
		keys := strings.Split(v, ",")
		tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[0], d); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", d); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	firstname := r.FormValue("firstname")
	// if firstname == "" {
	// }
	lastname := r.FormValue("lastname")
	// if lastname == "" {
	// }
	filters := r.FormValue("filters")
	// if filters == "" {
	// }

	iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
	p := dao.Patient{
		UUID:      mux.Vars(r)["uuid"],
		Firstname: firstname,
		Lastname:  lastname,
		Filters:   filters,
		Studies:   []dao.Study{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
	}

	if false {
		if err := iCtxV.DB.UpdatePatient(r.Context(), p); err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		keys := strings.Split(v, ",")
		tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[1], p); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", p); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	firstname := r.FormValue("firstname")
	// if firstname == "" {
	// }
	lastname := r.FormValue("lastname")
	// if lastname == "" {
	// }
	filters := r.FormValue("filters")
	// if filters == "" {
	// }

	iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
	p := dao.Patient{
		UUID:      uuid.NewString(),
		Firstname: firstname,
		Lastname:  lastname,
		Filters:   filters,
		Studies:   []dao.Study{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
	}

	// TODO: remove this when dynamodb is provisioned
	if false {
		if err := iCtxV.DB.AddPatient(r.Context(), p); err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		keys := strings.Split(v, ",")
		tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[1], p); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", p); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	}
}
