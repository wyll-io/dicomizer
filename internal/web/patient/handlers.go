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
	}

	if false {
		iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
		_, err := iCtxV.DB.GetPatient(r.Context(), dao.SearchPatientParams{
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
		Studies:   []dao.DCMImage{},
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

	d := map[string]interface{}{
		"UUID":      uuid.NewString(),
		"Firstname": r.FormValue("firstname"),
		"Lastname":  r.FormValue("lastname"),
		"Filters":   r.FormValue("filters"),
		"Studies":   []dao.DCMImage{},
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Time{},
		"DeletedAt": time.Time{},
	}
	var firstnameFormErr string
	var lastnameFormErr string
	var filtersFormErr string
	if d["Firstname"] == "" || len(d["Firstname"].(string)) < 2 {
		firstnameFormErr = "Le prénom est obligatoire et doit contenir au moins 2 caractères"
	}
	if d["Lastname"] == "" || len(d["Lastname"].(string)) < 2 {
		lastnameFormErr = "Le nom de famille est obligatoire et doit contenir au moins 2 caractères"
	}
	if d["Filters"] == "" || len(d["Filters"].(string)) < 2 {
		filtersFormErr = "Les filtres sont obligatoires et doivent contenir au moins 2 caractères"
	}

	if firstnameFormErr != "" || lastnameFormErr != "" || filtersFormErr != "" {
		d["Form"] = map[string]string{}
		if firstnameFormErr != "" {
			d["Form"].(map[string]string)["Firstname"] = firstnameFormErr
		}
		if lastnameFormErr != "" {
			d["Form"].(map[string]string)["Lastname"] = lastnameFormErr
		}
		if filtersFormErr != "" {
			d["Form"].(map[string]string)["Filters"] = filtersFormErr
		}

		if v := r.Header.Get("DICOMIZER-PARTIAL-400"); v != "" {
			keys := strings.Split(v, ",")
			tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusBadRequest)

			if err := tmpl.ExecuteTemplate(w, keys[1], d); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			webError.RedirectError(w, r, http.StatusInternalServerError, "Validation d'un formulaire sans passer par le GUI")
		}

		return
	}

	// TODO: remove this when dynamodb is provisioned
	if false {
		iCtxV := r.Context().Value(webContext.Internal).(webContext.InternalValues)
		if err := iCtxV.DB.AddPatient(r.Context(), dao.Patient{
			UUID:      d["UUID"].(string),
			Firstname: d["Firstname"].(string),
			Lastname:  d["Lastname"].(string),
			Filters:   d["Filters"].(string),
			Studies:   d["Studies"].([]dao.DCMImage),
			CreatedAt: d["CreatedAt"].(time.Time),
			UpdatedAt: d["UpdatedAt"].(time.Time),
			DeletedAt: d["DeletedAt"].(time.Time),
		}); err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		keys := strings.Split(v, ",")
		tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[1], d); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", d); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	}
}
