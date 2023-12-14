package error

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrInternalError = "Une erreur est survenue lors du traitement de votre requête"
	ErrParseForm     = "Une erreur est survenue lors du traitement du formulaire"
)

func RedirectError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	params := url.Values{}
	params.Add("code", strconv.Itoa(code))
	params.Add("message", msg)

	http.Redirect(w, r, fmt.Sprintf("/error?%s", params.Encode()), http.StatusMovedPermanently)
}
