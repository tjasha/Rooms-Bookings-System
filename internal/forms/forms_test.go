package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// function Valid() has a receiver, so we put receiver name in the test name
func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/something", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Errorf("Got invalid when it should be valid.")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/something", nil)
	form := New(r.PostForm)

	//test setting fields a, b and c as required
	form.Required("a", "b", "c")
	//form.Required()
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing.")
	}

	//setting values in the fields
	postedData := url.Values{}
	postedData.Set("a", "a")
	postedData.Set("b", "a")
	postedData.Set("c", "a")

	r, _ = http.NewRequest(http.MethodPost, "/something", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form  shows doesn't have required fields when it does.")
	}

}

// func (f *Form) Has(field string, r *http.Request) bool {
func TestForm_Has(t *testing.T) {
	//setting value in the field
	postedData := url.Values{}
	postedData.Add("a", "aaa")
	form := New(postedData)

	has := form.Has("a")
	if !has {
		t.Errorf("Field should have value but it doesn't")
	}

	// setting field without value
	postedData.Add("b", "")
	form = New(postedData)

	has = form.Has("b")
	if has {
		t.Errorf("Field has value but it shouldn't")
	}
}

//func /(f *Form) MinLength(field string, length int, r *http.Request) bool {

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "aaa")
	form := New(postedData)

	length := 3
	minLength := form.MinLength("a", length)
	if !minLength {
		t.Errorf("Field should have at least length %v but it doesn't", length)
	}

	// this is a test for Get() function
	isError := form.Errors.Get("a")
	if isError != "" {
		t.Errorf("Should not have an error but got one ")
	}

	postedData.Add("b", "bb")
	minLength = form.MinLength("b", length)
	if minLength {
		t.Errorf("Field length is meeting requirement")
	}

	minLength = form.MinLength("c", length)
	if form.Valid() {
		t.Errorf("Form shows min length for non-existing field")
	}

	// this is a test for Get() function
	isError = form.Errors.Get("c")
	if isError == "" {
		t.Errorf("Should have an error but didn't get one ")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "test@email.com")
	form := New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Errorf("Form showm invalid email when email is valid")
	}

	postedData.Add("notemail", "testemailcom")
	form = New(postedData)

	form.IsEmail("notemail")
	if form.Valid() {
		t.Errorf("Confirmed email when it's invalid")
	}

	form.IsEmail("someField")
	if form.Valid() {
		t.Errorf("Form shows emil for non-existing field")
	}
}
