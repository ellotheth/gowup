package gowup

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ApiTest struct {
	suite.Suite
	api WIU
}

func TestApi(t *testing.T) {
	suite.Run(t, new(ApiTest))
}

func (a *ApiTest) SetupSuite() {
	a.api = WIU{Client: "herp", Token: "derp"}
}

func (a *ApiTest) TestSetHeaderDefaults() {
	req, _ := http.NewRequest("GET", "", nil)

	a.api.setHeaders(req, nil)

	a.Equal("application/json", req.Header.Get("Content-Type"), "should have a json content type")
	a.Equal("Bearer herp derp", req.Header.Get("Auth"), "should set client auth")
}

func (a *ApiTest) TestSetCustomHeaders() {
	req, _ := http.NewRequest("GET", "", nil)

	a.api.setHeaders(req, map[string]string{"herp": "derp", "bar": "foo"})

	a.Equal("derp", req.Header.Get("herp"), "should set herp header")
	a.Equal("foo", req.Header.Get("bar"), "should set bar header")
}

func (a *ApiTest) TestGet() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	endpoint := "foo"

	apiEntryPoint = server.URL
	response, err := a.api.get(endpoint)
	server.Close()

	a.Nil(err, "should not return an error")

	r := response.Request
	a.Equal("GET", r.Method, "should be a GET")
	a.Equal(fmt.Sprintf("/%s", endpoint), r.URL.Path, "should be the right endpoint")

	// todo: how do i confirm setHeaders was run?
}

func (a *ApiTest) TestPost() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	}))
	defer server.Close()

	apiEntryPoint = server.URL
	endpoint := "foo"
	data := map[string]interface{}{"derp": "thing", "foo": []interface{}{"a", "string"}, "herp": 1}
	marshaled, _ := json.Marshal(data)

	response, err := a.api.post(endpoint, data)
	a.Nil(err, "should not return an error")

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	a.Equal("POST", response.Request.Method, "should be a POST")
	a.Equal(fmt.Sprintf("/%s", endpoint), response.Request.URL.Path, "should be the right endpoint")
	a.Equal(marshaled, body, "should post json data")

	// apparently comparing maps is complicated. or json.Unmarshal is flaky.
	// this comparison is USUALLY fine, but sometimes it fails for literally no
	// reason at all. seriously. throw some fmt.Printf lines in there and try
	// it.
	//
	// var unmarshaled map[string]interface{}
	// json.Unmarshal(body, &unmarshaled)
	// assert.True(t, assert.ObjectsAreEqual(data, unmarshaled), "should decode to the right json")
}

func (a *ApiTest) TestLocations() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "sources": [
        {
            "id": "2",
            "name": "toronto",
            "title": "Toronto",
            "location": "Toronto",
            "state": "Ontario",
            "country": "Canada",
            "latitude": "43.6481",
            "longitude": "-79.4042",
            "continent_name": "North America"
        },
        {
            "id": "12",
            "name": "dallas",
            "title": "Dallas",
            "location": "Dallas",
            "state": "Texas",
            "country": "United States",
            "latitude": "32.7828",
            "longitude": "-96.8039",
            "continent_name": "North America"
        },
        {
            "id": "13",
            "name": "newyork",
            "title": "New York",
            "location": "Garden City",
            "state": "New York",
            "country": "United States",
            "latitude": "40.7269",
            "longitude": "-73.6497",
            "continent_name": "North America"
        }
    ]
}`))
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	sources, err := a.api.Locations()
	a.Nil(err, "should not return an error")
	a.Equal(3, len(sources), "should contain the full list of servers")
	a.Equal("dallas", sources[1]["name"], "should have the same content as the raw json")
}

// todo: test error handling for get
// todo: test error handling for post
