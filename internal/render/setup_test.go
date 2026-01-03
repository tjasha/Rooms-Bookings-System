package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

// it gets called before any test get run, execute whatever we want and run tests just before exiting
func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	//change this to true when in production, using it to define encription
	testApp.InProduction = false

	//adding logging to test setup
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	//initiate session package
	session = scs.New()
	session.Lifetime = 24 * time.Hour              //i want session to persist for 24h
	session.Cookie.Persist = true                  // session will be stored in the cookie
	session.Cookie.SameSite = http.SameSiteLaxMode // strict about the sites that cookie is valid for
	session.Cookie.Secure = false                  // no need for encription

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// we're creating http response writer
// response interface has header, writeHeader amd write
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

// it can stay empty
func (tw *myWriter) WriteHeader(i int) {

}

// writer
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
