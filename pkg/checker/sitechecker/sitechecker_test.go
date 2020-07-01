package sitechecker_test

import (
	"gohealthcheck/pkg/checker"
	"gohealthcheck/pkg/checker/sitechecker"
	"testing"

	"golang.org/x/xerrors"
	. "gopkg.in/check.v1"
)

type sitesTestData struct {
	checkeeType checker.CheckeeType
	url         string
	outcome     bool
	err         error
}

const (
	wrongType checker.CheckeeType = 5
)

var (

	// Reference for how to test for request timeout error:
	// https://stackoverflow.com/questions/100841/artificially-create-a-connection-timeout-error
	testTableSites = []sitesTestData{
		{checkeeType: checker.TypeSite, url: "https://www.google.com/", outcome: true, err: nil},
		{checkeeType: wrongType, url: "https://stackoverflow.com/", outcome: false, err: sitechecker.ErrorWrongCheckeeType},
		{checkeeType: checker.TypeSite, url: "http://10.255.255.1/", outcome: false, err: nil},
	}
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestSiteChecker(c *C) {

	siteChecker := sitechecker.NewSiteChecker()

	for _, site := range testTableSites {
		success, err := siteChecker.Check(checker.Checkee{Type: site.checkeeType, URL: site.url})
		c.Assert(xerrors.Is(err, site.err), Equals, true)
		c.Assert(success, Equals, site.outcome)
	}
}
