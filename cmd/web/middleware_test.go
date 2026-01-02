package main

import (
	"fmt"
	"net/http"
	"testing"
)

// here we need to set up interfaces, so
func TestNoSurf(t *testing.T) {
	var myH myHandler //this is defined in setup_test

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Error(fmt.Printf("type is not http.Handler, but is %T", v))

	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler //this is defined in setup_test

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
	// do nothing
	default:
		t.Error(fmt.Printf("type is not http.Handler, but is %T", v))

	}
}
