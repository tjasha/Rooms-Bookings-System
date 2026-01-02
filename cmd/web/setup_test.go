package main

import (
	"net/http"
	"os"
	"testing"
)

// This is setting up something to run before the test, run the test, something after the test and exit
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
