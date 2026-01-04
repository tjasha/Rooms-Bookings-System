package main

import (
	"encoding/gob"
	"fmt"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/driver"
	"github.com/tjasha/Rooms-Bookings-System/internal/handlers"
	"github.com/tjasha/Rooms-Bookings-System/internal/helpers"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"github.com/tjasha/Rooms-Bookings-System/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

// we have to run it with "go run *.go" now

var app config.AppConfig // now we can also use it in routes
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	//we're closing connection here not in the run()
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	//we add something that actually serves
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	//we need to start a server
	err = srv.ListenAndServe()
	log.Panic(err)
	//log.Fatal(err)
}

func run() (*driver.DB, error) {

	//what can i put in the session - primitive types are already allowed
	// we want to store reservation object
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})

	//change this to true when in production, using it to define encription
	app.InProduction = false

	//Stdout = standart out = terminal window
	//prefix to all of this logs is INFO
	//adding date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//adding info log to the application
	app.InfoLog = infoLog

	// currently writing in the terminal window
	// prefix will be ERROR
	// adding date and time, but also more information about the error
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//adding error log to the application
	app.ErrorLog = errorLog

	//initiate session package
	session = scs.New()
	session.Lifetime = 24 * time.Hour              //i want session to persist for 24h
	session.Cookie.Persist = true                  // session will be stored in the cookie
	session.Cookie.SameSite = http.SameSiteLaxMode // strict about the sites that cookie is valid for
	session.Cookie.Secure = app.InProduction       //this makes session encripted. while using localhost should be false, but in production should be true

	app.Session = session

	//connect to database
	log.Println("Connectiong to DB")
	db, err := driver.ConnectSQL("host=localhost port=5445 dbname=booking user=tjasaspes password=")

	if err != nil {
		log.Fatal("Cannot connect to database! Dying..")
	}

	log.Println("Connectiong to DB")

	// in theory closing connection should be here, but if it is, it'll close connection just after opening
	// for this reason, we return *driver.DB from the run()
	//defer db.SQL.Close()

	//i want to create template cache here
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(fmt.Printf("cannot create template cache %v", err))
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	//this give render access to appConfig
	render.NewRenderer(&app)

	//create repository variable (app config and database connection)
	repo := handlers.NewRepo(&app, db)
	//create handlers and return variable back to handlers
	handlers.NewHandlers(repo)

	//this will populate app in helpers.go with app config
	helpers.NewHelpers(&app)

	return db, nil
}
