package authentication

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	webContext "github.com/wyll-io/dicomizer/internal/web/context"
	webError "github.com/wyll-io/dicomizer/internal/web/error"
)

func AccessGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(cookieKey)
		if err != nil && err != http.ErrNoCookie {
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
			return
		} else if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		if err := c.Valid(); err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    cookieKey,
				Expires: time.Unix(0, 0),
				Path:    "/",
				Value:   "",
			})

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		} else if c.Value == "" {
			http.SetCookie(w, &http.Cookie{
				Name:    cookieKey,
				Expires: time.Unix(0, 0),
				Path:    "/",
				Value:   "",
			})

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		t, err := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		switch {
		case t.Valid:
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), webContext.User, true)))
		case errors.Is(err, jwt.ErrTokenMalformed):
			http.SetCookie(w, &http.Cookie{
				Name:    cookieKey,
				Expires: time.Unix(0, 0),
				Path:    "/",
				Value:   "",
			})

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			http.SetCookie(w, &http.Cookie{
				Name:    cookieKey,
				Expires: time.Unix(0, 0),
				Path:    "/",
				Value:   "",
			})

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			http.SetCookie(w, &http.Cookie{
				Name:    cookieKey,
				Expires: time.Time{},
				Path:    "/",
				Value:   "",
			})

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		default:
			webError.RedirectError(w, r, http.StatusInternalServerError, err.Error())
		}
	})
}
