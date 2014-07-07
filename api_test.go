package gowup

import (
	"encoding/json"
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

func (a *ApiTest) TestParse() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := json.Marshal([]int{1, 2, 3})
		w.Write(raw)
	}))
	defer server.Close()
	response, _ := http.Get(server.URL)

	dest := make([]int, 3)
	err := a.api.parse(response, &dest)
	a.NoError(err, "should not return an error")
	a.Equal([]int{1, 2, 3}, dest, "should have the right type and value")
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
	a.Equal("/"+endpoint, r.URL.Path, "should be the right endpoint")

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
	a.Equal("/"+endpoint, response.Request.URL.Path, "should be the right endpoint")
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

	expected := Location{
		Name:      "newyork",
		Title:     "New York",
		City:      "Garden City",
		State:     "New York",
		Country:   "United States",
		Lat:       "40.7269",
		Lon:       "-73.6497",
		Continent: "North America",
	}

	sources, err := a.api.Locations()
	a.Nil(err, "should not return an error")
	a.Equal(3, len(sources), "should contain the full list of servers")
	a.Equal(expected, sources[2], "should have the same content as the raw json")
}

func (a *ApiTest) TestJobs() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		    "534419e98c3dcffa6170aeae": {
		        "url": "https://google.com",
		        "ip": "123.4.56.189",
		        "start_time": 1396972009,
		        "easy_time": "Tue, 08 Apr 2014 11:46:49 -0400",
		        "services": [
		            {
		                "city": "denver",
		                "server": "denver",
		                "checks": [
		                    "http",
		                    "ping",
		                    "trace",
		                    "fast",
		                    "dig"
		                ]
		            },
		            {
		                "city": "sydney",
		                "server": "sydney",
		                "checks": [
		                    "http",
		                    "ping",
		                    "trace",
		                    "fast",
		                    "dig"
		                ]
		            },
		            {
		                "city": "riga",
		                "server": "riga",
		                "checks": [
		                    "http",
		                    "ping",
		                    "trace",
		                    "fast",
		                    "dig"
		                ]
		            }
		        ]
		    }
		}`))
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	jobs, err := a.api.Jobs()
	a.NoError(err, "should not return an error")

	job, ok := jobs["534419e98c3dcffa6170aeae"]
	a.True(ok, "job should exist in jobs")

	a.Equal("https://google.com", job.Url.String(), "should unmarshal urls")
	a.Equal(1396972009, job.StartTime.Unix(), "should unmarshal time")
	a.Equal("123.4.56.189", job.Ip, "should unmarshal IP address")
	a.Equal("trace", job.Services[0].Tests[2], "should unmarshal services")
}

func (a *ApiTest) TestBadJobId() {
	_, err := a.api.Job("herpderp")
	a.Error(err, "should reject bad job ids")
}

func (a *ApiTest) TestJobWithSummaryOnly() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		    "request": {
		        "easy_time": "Sun, 29 Jun 2014 10:53:09 -0400",
		        "expiry": {
		            "sec": 1405263189,
		            "usec": 0
		        },
		        "ip": "107.4.56.245",
		        "start_time": 1404053589,
		        "url": "https://google.com"
		    }
		}`))
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	job, err := a.api.Job("aa")
	a.NoError(err, "should not return an error")

	a.Equal(1404053589, job.Summary.StartTime.Unix(), "should decode the start time")
	a.Equal(1405263189, job.Summary.ExpireTime.Unix(), "should decode the expiration time")
	a.Equal("107.4.56.245", job.Summary.Ip, "should decode the ip address")
	a.Equal("https://google.com", job.Summary.Url.String(), "should decode the url")
}

func (a *ApiTest) TestJobDetails() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		    "request": {
		        "easy_time": "Sun, 29 Jun 2014 10:53:09 -0400",
		        "expiry": {
		            "sec": 1405263189,
		            "usec": 0
		        },
		        "ip": "107.4.56.245",
		        "start_time": 1404053589,
		        "url": "https://google.com"
		    },
		    "response": {
		        "complete": {
		            "denver": {
		                "fast": {
		                    "raw": { "some": "random content" },
		                    "summary": { "some": "random content" }
		                }
		            }
		        },
		        "error": [],
		        "in_progress": []
		    }
		}`))
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	job, err := a.api.Job("aa")
	a.NoError(err, "should not return an error")
	a.Equal(
		map[string]interface{}{"some": "random content"},
		job.Details.Done["denver"]["fast"],
		"should decode deeply nested json",
	)
	a.Equal(JobDetail{}, job.Details.Error, "should decode empty details to nil")
}

func (a *ApiTest) TestEmptySubmission() {
	_, err := a.api.Submit(nil)
	a.Error(err, "should error on an empty submission")
	a.Equal("Nothing to submit", err.Error())
}

func (a *ApiTest) TestSubmissionPostFailure() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "derp", 400)
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	_, err := a.api.Submit(&JobRequest{})
	a.Error(err, "should throw an error if the server breaks")
	a.Contains(err.Error(), "invalid character")
}

func (a *ApiTest) TestJobRequestMarshaling() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var req JobRequest
		json.Unmarshal(body, &req)

		a.Equal("https://foo/bar?herp=derp", req.Url, "should encode the full url")
		a.Equal([]string{"foo", "bar"}, req.Tests, "should encode the tests array")
		a.Equal([]string{"herp", "derp"}, req.Locations, "should encode sources array")
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	req := JobRequest{
		Url:       "https://foo/bar?herp=derp",
		Tests:     []string{"foo", "bar"},
		Locations: []string{"herp", "derp"},
	}
	a.api.Submit(&req)
}

func (a *ApiTest) TestBadServerData() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal("NOPE")
		w.Write(data)
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	_, err := a.api.Submit(&JobRequest{})
	a.Error(err, "should throw an error if the data is the wrong format")
	a.Contains(err.Error(), "cannot unmarshal string")
}

func (a *ApiTest) TestMissingJobID() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(map[string]string{"message": "nope"})
		w.Write(data)
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	_, err := a.api.Submit(&JobRequest{})
	a.Error(err, "should throw an error if there's no job id")
	a.Equal("Submission failed", err.Error())
}

func (a *ApiTest) TestJobID() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(map[string]string{"jobID": "herpderp"})
		w.Write(data)
	}))
	defer server.Close()
	apiEntryPoint = server.URL

	id, err := a.api.Submit(&JobRequest{})
	a.NoError(err, "should not return an error")
	a.Equal("herpderp", id, "should match the job id sent")
}

// todo: test error handling for get
// todo: test error handling for post
