package models

import "github.com/tjasha/Rooms-Bookings-System/internal/forms"

// this struct holds any possible kind of data that we could send to the template
// holds data set from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} //we're not sure what kind of data that may be, so usinf interface
	CSRFToken string                 //this would be forms
	Flash     string                 // messaged
	Warning   string                 //warnings
	Error     string
	Form      *forms.Form //allows that i can display empty form before user submit it
}
