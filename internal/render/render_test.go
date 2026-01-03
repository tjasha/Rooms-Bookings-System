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

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	//we need request
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	//we need response (created in setup_test)
	var ww myWriter

	//Test1 : existing template
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	//Test2: non-existing template
	err = RenderTemplate(&ww, r, "non-existing.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template doesn't exist")
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

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
