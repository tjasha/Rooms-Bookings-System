package forms

import (
	"net/http"
	"net/url"
)

// information about the form
type Form struct {
	url.Values
	Errors errors
}

// initialize a form struct - it's expected to be empty in the beginning
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks that there's a value in the required field
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}
