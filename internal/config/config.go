package config

import (
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

//we should not allow to import anything in this package
//we should import this package everywhere that we need

// we can put here anything that we need "globally" in our application
type AppConfig struct {
	UseCache      bool //we use this to allow us to test in dev mode
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger // to handle errors easier
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData // creating channel that is available al the time from anywhere in the application.
}
