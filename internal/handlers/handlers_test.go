package handlers

import (
	"context"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	//{"home", "/", "GET", []postData{}, http.StatusOK},
	//{"about", "/about", "GET", []postData{}, http.StatusOK},
	//{"general-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	//{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	//{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	//{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	//{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	//
	//{"post-search-availability", "/search-availability", "POST", []postData{
	//	{key: "start", value: "02-01-2026"},
	//	{key: "end", value: "03-01-2026"},
	//}, http.StatusOK},
	//{"post-search-availability-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "02-01-2026"},
	//	{key: "end", value: "03-01-2026"},
	//}, http.StatusOK},
	//{"post-make-reservation", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "John"},
	//	{key: "last_name", value: "Smith"},
	//	{key: "email", value: "john.smitg@email.com"},
	//	{key: "phone", value: "0248-0695-5465"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	//we create server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == http.MethodGet {
			//this is creating the client
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {
			//we populate to value data that we are expecting from POST
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	//we get a request to put this info in the session
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//new recorder - it's simulating what we get from request
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	// we're turning handler in handler function, because we can't call it directly
	handler := http.HandlerFunc(Repo.Reservation)

	// we add response recorder and request that we built
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func getCtx(req *http.Request) context.Context {
	//we're getting header X-Session that we need later!
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
