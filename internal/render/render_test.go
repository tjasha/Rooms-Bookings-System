package render

import (
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	//creating request
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// what we want to test
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	//test results
	if result.Flash != "123" {
		t.Error("flash value 123 not found in session")
	}

}

// i need to add context to the request
func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	//this is making it active session
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	//adding session in request
	r = r.WithContext(ctx)

	return r, nil
}
