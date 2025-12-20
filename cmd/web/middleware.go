package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// All middlewares need to have a parameter usually called next type http.Handler
// all middleware needs to return http.Handler

//creating noSerf token
// adds CSRS protection to all POST requests
func NoSerf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	//have to creat a cookie with som values, valid to all sites
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// automatically loads a session and communicate cookie to and from middleware
// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}