package csvreader_test

import (
	"gohealthcheck/reader/csvreader"
	"io/ioutil"
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
`
	csvTestSampleData = []string{
		"https://www.google.com/",
		"https://stackoverflow.com/",
		"http://ihatemylife.org/"}
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

func (s *TestSuite) TestCSVReader(c *C) {
	csvReader, err := csvreader.NewCSVReader(filepath.Join(s.dir, fileName))
	c.Assert(err, IsNil)

	sites, err := csvReader.ReadSite()
	c.Assert(err, IsNil)

	dataIndex := 0
	for sites.Next() {
		c.Assert(sites.Site().URL, Equals, csvTestSampleData[dataIndex])
		dataIndex++
	}

	err = sites.Close()
	c.Assert(err, IsNil)

}
