package main

import (
	"github.com/tjasha/Rooms-Bookings-System/internal/helpers"
	"net/http"

	"github.com/justinas/nosurf"
)

// All middlewares need to have a parameter usually called next type http.Handler
// all middleware needs to return http.Handler

// creating noSurf token
// adds CSRS protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	//have to creat a cookie with som values, valid to all sites
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// automatically loads a session and communicate cookie to and from middleware
// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// i need to have access to request here
// i can apply to any routs that i want to protect
func Auth(next http.Handler) http.Handler {

	// to get access to request, i call return on HandlerFunc, that has anonymous function, that returns request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// i put request here
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// this is just passing flow further
		next.ServeHTTP(w, r)
	})
}
