package guppy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetHeaderDefaults(t *testing.T) {
	api := WIU{Client: "herp", Token: "derp"}
	req, _ := http.NewRequest("GET", "", nil)

	api.setHeaders(req, nil)

	assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "should have a json content type")
	assert.Equal(t, "Bearer herp derp", req.Header.Get("Auth"), "should set client auth")
}

func TestSetCustomHeaders(t *testing.T) {
	api := WIU{}
	req, _ := http.NewRequest("GET", "", nil)

	api.setHeaders(req, map[string]string{"herp": "derp", "bar": "foo"})

	assert.Equal(t, "derp", req.Header.Get("herp"), "should set herp header")
	assert.Equal(t, "foo", req.Header.Get("bar"), "should set bar header")
}

func TestGet(t *testing.T) {
	endpoint, client, token := "foo", "herp", "derp"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	apiEntryPoint = server.URL
	api := WIU{Client: client, Token: token}
	response, err := api.get(endpoint)

	assert.Nil(t, err, "should not return an error")

	r := response.Request
	assert.Equal(t, "GET", r.Method, "should be a GET")
	assert.Equal(t, fmt.Sprintf("/%s", endpoint), r.URL.Path, "should be the right endpoint")
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "should have a json content type")
	assert.Equal(t, fmt.Sprintf("Bearer %s %s", client, token), r.Header.Get("Auth"), "should set client auth")
}
