package http_test

import (
	"encoding/json"
	"fmt"
	"gohealthcheck/pkg/reporter"
	"net/http"
	"net/http/httptest"
	"testing"

	h "gohealthcheck/pkg/reporter/http"

	. "gopkg.in/check.v1"
)

var (
	stats = reporter.Statistics{
		TotalWebSites: 3,
		Success:       2,
		Failure:       1,
		TotalTime:     1,
	}
	token = "token1234"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpTest(c *C) {
}

func (s *TestSuite) TestHTTPReporter(c *C) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		c.Assert(len(r.Header["Authorization"]), Equals, 1)
		c.Assert(r.Header["Authorization"][0], Equals, fmt.Sprintf("Bearer %s", token))

		req := reporter.Statistics{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		c.Assert(err, IsNil)

		c.Assert(req.TotalWebSites, Equals, stats.TotalWebSites)
		c.Assert(req.Success, Equals, stats.Success)
		c.Assert(req.Failure, Equals, stats.Failure)
		c.Assert(req.TotalTime, Equals, stats.TotalTime)

	}))
	defer ts.Close()

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	reporter := h.NewHTTPReporter(ts.URL, headers)
	reporter.Report(stats)
}
