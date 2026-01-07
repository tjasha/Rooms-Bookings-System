package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/driver"
	"github.com/tjasha/Rooms-Bookings-System/internal/forms"
	"github.com/tjasha/Rooms-Bookings-System/internal/helpers"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"github.com/tjasha/Rooms-Bookings-System/internal/render"
	"github.com/tjasha/Rooms-Bookings-System/internal/repository"
	"github.com/tjasha/Rooms-Bookings-System/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

//Repository pattern:

// repository used by handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// create new repository
// also adding database connection
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostresRepo(db.SQL, a),
	}
}

// sets repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// we add receiver(m *Repository) to all handlers - this give them an access to the application variables
// Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// to test sessions we save ip address of the user coming to the home page in the session
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders reservation page and display the form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	//we pull reservation out of the session (it has start and end date and roomID
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("reservation session not found"))
		return
	}

	// ret information of the room from DB
	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// save room name into session
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	// we need to change dates from time.Time to displayable date
	sd := res.StartDate.Format("02-01-2006")
	ed := res.EndDate.Format("02-01-2006")
	//adding dates in the map
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap, //adding dates
	})
}

// PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("reservation session not found"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		//we can use centralised errors
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	//get data that user inputed
	form := forms.New(r.PostForm)

	//form.Has("first_name", r) --> it's not necessary anymore
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		//saved data that user wrote
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// write information in the database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// we put reservation object to the session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	//when you receive POST request, by good practice, you should direct the user through http redirect
	//StatusSeeOther is a standard
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the general quarters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the majors suite page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders search availability
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders search availability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, room := range rooms {
		m.App.InfoLog.Println("Room: ", room.ID, room.RoomName)
	}

	if len(rooms) == 0 {
		//no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	// we create new reservation that will be stored in the session.
	// i need to store start date, end date and roomID (will be populated later
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON test
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "02-01-2006"
	startDate, err := time.Parse(layout, sd)
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, _ := m.DB.SearchAvailabilityForDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK:      available,
		Message: "",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders contact
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//asserting reservation into Reservation object
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		//if user gets to reservation summary page without reservation, this error will be sent and user will be redirected
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//when we don't need reservation object anymore, we can remove it from session
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("02-01-2006")
	ed := reservation.EndDate.Format("02-01-2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	//we want to read the id from the link and store it as a roomID (id from the url)
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// i need to get reservation from the session, update with roomID and store it back to the session
	//we get session value by name (it's interface) and cast it to the model.Reservation
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	// save roomID
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	// we redirect to the reservation page
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
