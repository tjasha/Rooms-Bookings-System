package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"github.com/tjasha/Rooms-Bookings-System/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	//change this to true when in production, using it to define encription
	app.InProduction = false

	//initiate session package
	session = scs.New()
	session.Lifetime = 24 * time.Hour              //i want session to persist for 24h
	session.Cookie.Persist = true                  // session will be stored in the cookie
	session.Cookie.SameSite = http.SameSiteLaxMode // strict about the sites that cookie is valid for
	session.Cookie.Secure = app.InProduction       //this makes session encripted. while using localhost should be false, but in production should be true

	app.Session = session

	//i want to create template cache here
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal(fmt.Printf("cannot create template cache %v", err))
	}

	app.TemplateCache = tc
	//we don't want to create new templates on every request
	app.UseCache = true

	//create repository variable
	repo := NewRepo(&app)
	//create handlers and return variable back to handlers
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	//we can easily add some middleware that is available in chi
	//we are adding middleware that is recovering application if there was fatal error
	mux.Use(middleware.Recoverer)
	// we are calling middleware here - this is security protection! don't remove it!
	// mux.Use(NoSurf) --> to test, we don't want to look for crsf token
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	//where all our images are saved
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{} //it's creating empty map

	// i want to add automatically all available templates that exist
	// they should be added in order

	// i want to first add all *page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) //we just look for all files with this pattern
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *page.tmpl that we found before
	for _, page := range pages { //page is full path to the template
		name := filepath.Base(page)                                     //returns last element of the path = name of the file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) //(ts = template set) we  parse this file (page) and store in the template (name)
		if err != nil {
			return myCache, err
		}

		//now we look for all layouts - we use the same syntax as for the pages
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		//checking how many elements we have
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates)) // check if any of the pages needs layout inside of them to be rendered. if yes, it adds it to the ts
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts //adding template set to cache map
	}

	return myCache, nil
}
