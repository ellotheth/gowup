package guppy

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

func (apiSuite *ApiTest) SetupSuite() {
	apiSuite.api = WIU{Client: "herp", Token: "derp"}
}

func (apiSuite *ApiTest) TestSetHeaderDefaults() {
	req, _ := http.NewRequest("GET", "", nil)

	apiSuite.api.setHeaders(req, nil)

	apiSuite.Equal("application/json", req.Header.Get("Content-Type"), "should have a json content type")
	apiSuite.Equal("Bearer herp derp", req.Header.Get("Auth"), "should set client auth")
}

func (apiSuite *ApiTest) TestSetCustomHeaders() {
	req, _ := http.NewRequest("GET", "", nil)

	apiSuite.api.setHeaders(req, map[string]string{"herp": "derp", "bar": "foo"})

	apiSuite.Equal("derp", req.Header.Get("herp"), "should set herp header")
	apiSuite.Equal("foo", req.Header.Get("bar"), "should set bar header")
}

func (apiSuite *ApiTest) TestGet() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	endpoint := "foo"

	apiEntryPoint = server.URL
	response, err := apiSuite.api.get(endpoint)
	server.Close()

	apiSuite.Nil(err, "should not return an error")

	r := response.Request
	apiSuite.Equal("GET", r.Method, "should be a GET")
	apiSuite.Equal(fmt.Sprintf("/%s", endpoint), r.URL.Path, "should be the right endpoint")

	// todo: how do i confirm setHeaders was run?
}

func (apiSuite *ApiTest) TestPost() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	}))
	defer server.Close()

	apiEntryPoint = server.URL
	endpoint := "foo"
	data := map[string]interface{}{"derp": "thing", "foo": []interface{}{"a", "string"}, "herp": 1}
	marshaled, _ := json.Marshal(data)

	response, err := apiSuite.api.post(endpoint, data)
	apiSuite.Nil(err, "should not return an error")

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	apiSuite.Equal("POST", response.Request.Method, "should be a POST")
	apiSuite.Equal(fmt.Sprintf("/%s", endpoint), response.Request.URL.Path, "should be the right endpoint")
	apiSuite.Equal(marshaled, body, "should post json data")

	// apparently comparing maps is complicated. or json.Unmarshal is flaky.
	// this comparison is USUALLY fine, but sometimes it fails for literally no
	// reason at all. seriously. throw some fmt.Printf lines in there and try
	// it.
	//
	// var unmarshaled map[string]interface{}
	// json.Unmarshal(body, &unmarshaled)
	// assert.True(t, assert.ObjectsAreEqual(data, unmarshaled), "should decode to the right json")
}

func (apiSuite *ApiTest) TestLocations() {
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

	sources, err := apiSuite.api.Locations()
	apiSuite.Nil(err, "should not return an error")
	apiSuite.Equal(3, len(sources), "should contain the full list of servers")
	apiSuite.Equal("dallas", sources[1]["name"], "should have the same content as the raw json")
}

// todo: test error handling for get
// todo: test error handling for post
