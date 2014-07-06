package gowup

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
	"time"
)

type JobSummaryTest struct {
	suite.Suite
}

func TestJobSummary(t *testing.T) {
	suite.Run(t, new(JobSummaryTest))
}

func (j *JobSummaryTest) TestUrlUnmarshaler() {
	actual, urlstring := Url{}, "https://google.com/asdf"
	data, _ := json.Marshal(urlstring)

	parsed, _ := url.Parse(urlstring)
	expected := Url{URL: parsed}

	json.Unmarshal(data, &actual)
	j.Equal(expected, actual, "should unmarshal urls to a Url type")
	j.Equal(actual.Path, "/asdf", "should have the correct url content")
}

func (j *JobSummaryTest) TestUrlUnmarshalerInSummary() {
	job, urlstring := JobSummary{}, "https://herp.us/derp?foo=bar"
	data, _ := json.Marshal(map[string]string{"url": urlstring})

	parsed, _ := url.Parse(urlstring)
	expected := Url{URL: parsed}

	json.Unmarshal(data, &job)
	j.Equal(expected, job.Url, "should unmarshal url fields to Urls")
}

func (j *JobSummaryTest) TestTimeUnmarshaler() {
	job, start := JobSummary{}, Time{Time: time.Unix(1396972009, 0)}
	data, _ := json.Marshal(map[string]int64{"start_time": start.Unix()})

	json.Unmarshal(data, &job)
	j.Equal(start, job.StartTime, "should unmarshal time fields correctly")
}

func (j *JobSummaryTest) TestServicesUnmarshaler() {
	data := []byte(`{ "services": [
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
	]}`)

	job := JobSummary{}
	json.Unmarshal(data, &job)

	j.Equal("http", job.Services[2].Tests[0], "should unmarshal the nested slices")
	j.Equal("sydney", job.Services[1].Server, "should unmarshal the strings")
}
