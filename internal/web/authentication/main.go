package authentication

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

const cookieKey = "DICOMIZER_SESSION"

func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		login(w, r)
	case http.MethodPost:
		authenticate(w, r)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieKey,
		Expires: time.Unix(0, 0),
		Path:    "/",
		Value:   "",
	})

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := r.Context().Value(webContext.Templates).(webContext.TemplatesValues)["login"]
	if err := tmpl.ExecuteTemplate(w, "base.html", struct {
		Title    string
		LoggedIn bool
	}{LoggedIn: false, Title: "Connexion - Dicomizer"}); err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		webError.RedirectError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if r.Form.Get("password") != os.Getenv("ADMIN_PASSWORD") {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"iss": "dicomizer-web",
		"aud": "dicomizer-gui",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	token, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
