package handlers

import (
	"net/http"

	"github.com/tjasha/Rooms-Bookings-System/pkg/config"
	"github.com/tjasha/Rooms-Bookings-System/pkg/models"
	"github.com/tjasha/Rooms-Bookings-System/pkg/render"
)

//Repository pattern:

//repository used by handlers
var Repo *Repository

//repository type
type Repository struct {
	App *config.AppConfig
}

// create new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// sets repository for handlers
func NewHandlers(r *Repository){
	Repo = r
}


//we add receiver(m *Repository) to all handlers - this give them an access to the application variables
//Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// to test sessions we save ip address of the user coming to the home page in the session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}) 
} 

//About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	//some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	// testing session, we read saved ip from the session in about page
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, //we add map that we created to the template
	}) 

} 
