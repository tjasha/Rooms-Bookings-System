package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tjasha/Rooms-Bookings-System/pkg/config"
	"github.com/tjasha/Rooms-Bookings-System/pkg/handlers"
	"github.com/tjasha/Rooms-Bookings-System/pkg/render"
)

const portNumber = ":8080"

// we have to run it with "go run *.go" now

var app config.AppConfig // now we can also use it in routes
var session *scs.SessionManager


func main() {

	//change this to true when in production, using it to define encription 
	app.InProduction = false

	//initiate session package
	session = scs.New()
	session.Lifetime = 24 * time.Hour //i want session to persist for 24h
	session.Cookie.Persist = true // session will be stored in the cookie
	session.Cookie.SameSite = http.SameSiteLaxMode // strict about the sites that cookie is valid for 
	session.Cookie.Secure = app.InProduction //this makes session encripted. while using localhost should be false, but in production should be true 

	app.Session = session

	//i want to create template cache here
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false


	//this give render access to appConfig
	render.NewTemplates(&app)
	


	//create repository variable
	repo := handlers.NewRepo(&app)
	//create handlers and return variable back to handlers
	handlers.NewHandlers(repo)


	// we added this into routes
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil) //we specify what to listen, in this case localhost on port 8080

	//we add something that actually serves 
	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	//we need to start a server
	err = srv.ListenAndServe()
	log.Fatal(err)
}
