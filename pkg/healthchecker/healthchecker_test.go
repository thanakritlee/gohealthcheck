package healthchecker_test

import (
	"encoding/json"
	"fmt"
	"gohealthcheck/pkg/healthchecker"
	"gohealthcheck/pkg/reporter"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

var (
	fileName = "input.csv"

	csvTestSample = `https://www.google.com/
https://stackoverflow.com/
http://ihatemylife.org/
http://10.255.255.1/
`

	token = "sampletoken"

	totalWebSites   = 4
	numberOfSuccess = 3
	numberOfFailure = 1
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	dir string
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	outFile := filepath.Join(s.dir, fileName)
	err := ioutil.WriteFile(outFile, []byte(csvTestSample), os.ModePerm)
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestHealthChecker(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		c.Assert(len(r.Header["Authorization"]), Equals, 1)
		c.Assert(r.Header["Authorization"][0], Equals, fmt.Sprintf("Bearer %s", token))

		req := reporter.Statistics{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		c.Assert(err, IsNil)

		c.Assert(req.TotalWebSites, Equals, totalWebSites)
		c.Assert(req.Success, Equals, numberOfSuccess)
		c.Assert(req.Failure, Equals, numberOfFailure)

	}))

	input := healthchecker.Input{
		InputType:   healthchecker.TypeCSV,
		CSVFileName: filepath.Join(s.dir, fileName),
	}

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	report := []healthchecker.Report{
		{
			ReportType: healthchecker.TypeHTTP,
			URL:        ts.URL,
			Headers:    headers,
		},
		{
			ReportType: healthchecker.TypeConsole,
		},
	}

	checker, err := healthchecker.NewHealthChecker(&input, &report)
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckHealth()
	if err != nil {
		log.Fatal(err)
	}

	err = checker.ReportResult()
	if err != nil {
		log.Fatal(err)
	}

}
