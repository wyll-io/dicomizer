package patient

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router) {
	r.HandleFunc("/patient/search", search).Methods("POST")

	r.HandleFunc("/patient", create).Methods("POST")
	r.HandleFunc("/patient/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			delete(w, r)
		case http.MethodPut:
			update(w, r)
		}
	}).Methods("PUT", "DELETE")
}
