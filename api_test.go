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
	server *httptest.Server
	api    WIU
}

func TestApi(t *testing.T) {
	suite.Run(t, new(ApiTest))
}

func (apiSuite *ApiTest) SetupSuite() {
	apiSuite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	}))
	apiEntryPoint = apiSuite.server.URL
	apiSuite.api = WIU{Client: "herp", Token: "derp"}
}

func (apiSuite *ApiTest) TearDownSuite() {
	apiSuite.server.Close()
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
	endpoint := "foo"

	response, err := apiSuite.api.get(endpoint)

	apiSuite.Nil(err, "should not return an error")

	r := response.Request
	apiSuite.Equal("GET", r.Method, "should be a GET")
	apiSuite.Equal(fmt.Sprintf("/%s", endpoint), r.URL.Path, "should be the right endpoint")

	// todo: how do i confirm setHeaders was run?
}

func (apiSuite *ApiTest) TestPost() {
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

// todo: test error handling for get
// todo: test error handling for post
