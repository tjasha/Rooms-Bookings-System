package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/driver"
	"github.com/tjasha/Rooms-Bookings-System/internal/forms"
	"github.com/tjasha/Rooms-Bookings-System/internal/helpers"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"github.com/tjasha/Rooms-Bookings-System/internal/render"
	"github.com/tjasha/Rooms-Bookings-System/internal/repository"
	"github.com/tjasha/Rooms-Bookings-System/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// only for testing
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// get information of the room from DB
	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "error getting reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	//repopulate the form
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "Validation did not pass", http.StatusSeeOther)
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
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into the database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notifications

	// notification to the guest
	htmlMessage := fmt.Sprintf(`
		<strong>Reservation confirmation</strong><br>
		Dear %s %s, <br>
		This is confirming your reservation in %s from %s to %s.
`, reservation.FirstName, reservation.LastName, reservation.Room.RoomName,
		reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@gh.ds",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	//send message to the channel
	m.App.MailChan <- msg

	//notification to the owner
	htmlMessage = fmt.Sprintf(`
		<strong>Reservation confirmation</strong><br>
		This is confirming reservation for %s %s in %s from %s to %s.
`, reservation.FirstName, reservation.LastName, reservation.Room.RoomName,
		reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	msg = models.MailData{
		To:       "owner@gh.ds",
		From:     "me@gh.ds",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	//send message to the channel
	m.App.MailChan <- msg

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
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON test
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	//added later for testing purposes. we need to parse the form.
	err := r.ParseForm()
	if err != nil {

		// can't parse form, so return appropriate json, not error message
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "02-01-2006"
	startDate, err := time.Parse(layout, sd)
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse a date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse roomID")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	available, err := m.DB.SearchAvailabilityForDatesByRoomID(startDate, endDate, roomID)

	//added later for testing purposes
	if err != nil {
		// database return appropriate json, not error message
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	//removing this error because it's causing to go in circular loop
	//out, err := json.MarshalIndent(resp, "", "     ")
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders contact
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
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

// ChooseRoom choose available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	//we want to read the id from the link and store it as a roomID (id from the url)
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	helpers.ServerError(w, err)
	//	return
	//}
	//changed to test it easier
	//split URL by // and grab 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//In your test for ChooseRoom, you will want to set the URL on your request as follows:
	//req.RequestURI = "/choose-room/1"

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

// BookRoom takes URL parameters, builds sessional variable and take user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// we have values stored in the URL - id, s, e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, sd)
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	//prevent session fixation attack
	// we should use this always whit login and logout
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	// we want to send user back to the form
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	// try to authenticate the user
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		// if there is an error,i want to send user back to the log in form
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// we authenticate a user by saving the ID that we got in the session
	m.App.Session.Put(r.Context(), "user_id", id)
	// i want to send user to home page after authentication
	m.App.Session.Put(r.Context(), "flash", "Log in successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	// renew session token
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard renders the admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminAllReservations renders all new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations renders all reservations in admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows the reservation in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	// i need to know where the reservation came from (all or new page) and id of the reservation
	// get url from request, splitting it by /
	exploded := strings.Split(r.RequestURI, "/")
	//"/admin/reservation/all/5" -> id is in position 4
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	//we're adding src in the stringMap into the template
	stringMap := make(map[string]string)
	stringMap["src"] = src

	//get reservation from DB
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil), //to display it in form in admin-reservations-show.page.tmpl - to render the page
	})
}

// AdminPostShowReservation handles posting of reservation form
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get url from request, splitting it by /
	exploded := strings.Split(r.RequestURI, "/")
	//"/admin/reservation/all/5" -> id is in position 4
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	//we're adding src in the stringMap into the template
	stringMap := make(map[string]string)
	stringMap["src"] = src

	//get reservation from DB
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// in production you should add error checking here
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminProcessReservation marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	// update as processed - value 1
	_ = m.DB.UpdateProcessedForReservation(id, 1)
	// notify user that it's processed
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	// redirect to the same page as they came from
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminProcessReservation marks a reservation as processed
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	// delete reservation with id
	_ = m.DB.DeleteReservation(id)
	// notify user that it's deleted
	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	// redirect to the same page as they came from
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminReservationsCalendar display the reservation calendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	// assume that there's no month or year specified --> we want to show  today
	now := time.Now()

	// we want to have url with parameters /?y=2006&m=3
	// if the parameters exist, we want to reset now to this time
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	// how many years, months and days i want to add
	next := now.AddDate(0, 1, 0)
	// we don't substract, we use -value
	last := now.AddDate(0, -1, 0)

	// we need to format dates now --> golang formating
	// i want to get month
	nextMonth := next.Format("01")
	// year next month
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_year"] = now.Format("2006")

	// get teh first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1) // add one month and remove a day

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day() // number of days in the month

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
