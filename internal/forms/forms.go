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

// Valid returns true if no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
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
		f.Errors.Add(field, "Field "+field+" cannot be blank")
		return false
	}
	return true
}
