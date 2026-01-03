package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
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

// We can pass one or many things in this function
// Check for Required fields
func (f *Form) Required(fields ...string) {
	//it loops through all fields in the form
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks that there's a value in the required field
func (f *Form) Has(field string) bool {
	// i need to check the value from the receiver here, not request
	x := f.Get(field)
	if x == "" {
		// remove error to make it usable also for checkboxes
		//f.Errors.Add(field, "Field "+field+" cannot be blank")
		return false
	}
	return true
}

// MinLength check for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail is validating the email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
