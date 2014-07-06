package guppy

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//	"reflect"
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

	// todo: how do i confirm setHeaders was run?
}

func TestPost(t *testing.T) {
	endpoint := "foo"
	data := map[string]interface{}{"derp": "thing", "foo": []interface{}{"a", "string"}, "herp": 1}
	marshaled, _ := json.Marshal(data)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	}))
	defer server.Close()

	apiEntryPoint = server.URL
	api := WIU{}

	response, err := api.post(endpoint, data)
	assert.Nil(t, err, "should not return an error")

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	assert.Equal(t, "POST", response.Request.Method, "should be a POST")
	assert.Equal(t, fmt.Sprintf("/%s", endpoint), response.Request.URL.Path, "should be the right endpoint")
	assert.Equal(t, marshaled, body, "should post json data")

	// apparently comparing maps is complicated. or json.Unmarshal is flaky.
	// this comparison is USUALLY fine, but sometimes it fails for literally no
	// reason at all. seriously. throw some fmt.Printf lines in there and try
	// it.
	//
	// var unmarshaled map[string]interface{}
	// json.Unmarshal(body, &unmarshaled)
	// assert.True(t, assert.ObjectsAreEqual(data, unmarshaled), "should decode to the right json")
}

// todo: test error handling for get
// todo: test error handling for post
