package patient

import (
	"fmt"
	"net/http"
	"strings"

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

	ctx := r.Context()
	iCtxV := ctx.Value(webContext.Internal).(webContext.InternalValues)
	if err := iCtxV.DB.DeletePatient(ctx, fmt.Sprintf("PATIENT#%s", mux.Vars(r)["pk"])); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := iCtxV.S3.DeletePatientFiles(ctx, fmt.Sprintf("%s/%s", iCtxV.Center, mux.Vars(r)["pk"])); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	search := r.FormValue("search")
	var patients []dao.PatientInfo
	var err error
	ctx := r.Context()
	iCtxV := ctx.Value(webContext.Internal).(webContext.InternalValues)

	if search == "" {
		patients, err = iCtxV.DB.GetPatientsInfo(ctx)
		if err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		patients, err = iCtxV.DB.SearchPatientInfo(ctx, search)
		if err != nil {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		d := map[string]interface{}{
			"Patients": patients,
		}
		keys := strings.Split(v, ",")
		tmpl := ctx.Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
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

func update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	patient := dao.PatientInfo{
		PK:       fmt.Sprintf("PATIENT#%s", mux.Vars(r)["pk"]),
		SK:       "INFO#0",
		Fullname: r.FormValue("fullname"),
		Filters:  r.FormValue("filters"),
	}
	ctx := r.Context()
	iCtxV := ctx.Value(webContext.Internal).(webContext.InternalValues)

	err := iCtxV.DB.UpdatePatientInfo(ctx, patient.PK, &patient)
	if err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		patient.PK = strings.Replace(patient.PK, "PATIENT#", "", 1)
		keys := strings.Split(v, ",")
		tmpl := ctx.Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[1], patient); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", patient); err != nil {
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

	patient := dao.PatientInfo{
		Fullname: r.FormValue("fullname"),
		Filters:  r.FormValue("filters"),
	}
	formErrors := map[string]string{
		"Fullname": "",
		"Filters":  "",
	}
	ctx := r.Context()

	if patient.Fullname == "" || len(patient.Fullname) < 2 {
		formErrors["Fullname"] = "Le prénom est obligatoire et doit contenir au moins 2 caractères"
	}
	if patient.Filters == "" || len(patient.Filters) < 2 {
		formErrors["Filters"] = "Les filtres sont obligatoires et doivent contenir au moins 2 caractères"
	}

	if formErrors["Fullname"] != "" || formErrors["Filters"] != "" {
		if v := r.Header.Get("DICOMIZER-PARTIAL-400"); v != "" {
			keys := strings.Split(v, ",")
			tmpl := ctx.Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusBadRequest)

			if err := tmpl.ExecuteTemplate(w, keys[1], map[string]interface{}{
				"Errors":   formErrors,
				"Fullname": patient.Fullname,
				"Filters":  patient.Filters,
			}); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			webError.RedirectError(w, r, http.StatusInternalServerError, "Validation d'un formulaire sans passer par le GUI")
		}

		return
	}

	iCtxV := ctx.Value(webContext.Internal).(webContext.InternalValues)
	if err := iCtxV.DB.AddPatientInfo(ctx, &patient); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if v := r.Header.Get("DICOMIZER-PARTIAL"); v != "" {
		patient.PK = strings.Replace(patient.PK, "PATIENT#", "", 1)

		keys := strings.Split(v, ",")
		tmpl := ctx.Value(webContext.Templates).(webContext.TemplatesValues)[keys[0]]
		if len(keys) == 2 {
			if err := tmpl.ExecuteTemplate(w, keys[1], patient); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tmpl.ExecuteTemplate(w, "base.html", patient); err != nil {
				webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	} else {
		webError.RedirectError(w, r, http.StatusInternalServerError, "Validation d'un formulaire sans passer par le GUI")
	}
}
