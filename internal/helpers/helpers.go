package helpers

import (
	"fmt"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"net/http"
	"runtime/debug"
)

// we will add all helpers that will assist us with logging here

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// errors on client side
func ClientError(w http.ResponseWriter, status int) {
	//we can put info to the info log
	app.InfoLog.Println("Client error with star=tus of", status)
	// we want to return something to the user
	http.Error(w, http.StatusText(status), status)
}

// something is going wrong with the server
func ServerError(w http.ResponseWriter, err error) {
	// we want as much info about the error as possible
	// we print the error message and error stacktrace
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	// we want to return something to the user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
