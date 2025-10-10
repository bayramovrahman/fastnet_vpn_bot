package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/helpers"
	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ExtendedSessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session.Exists(r.Context(), "remember_me") {
			rememberMe := session.GetBool(r.Context(), "remember_me")
			if rememberMe {
				cookie := &http.Cookie{
					Name:     session.Cookie.Name,
					Path:     session.Cookie.Path,
					Domain:   session.Cookie.Domain,
					Secure:   session.Cookie.Secure,
					HttpOnly: session.Cookie.HttpOnly,
					SameSite: session.Cookie.SameSite,
					MaxAge:   int((7 * 24 * time.Hour).Seconds()),
				}

				token := session.Token(r.Context())
				cookie.Value = token

				http.SetCookie(w, cookie)
			}
		}

		next.ServeHTTP(w, r)
	})
}
