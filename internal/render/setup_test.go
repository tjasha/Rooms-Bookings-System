package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
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
