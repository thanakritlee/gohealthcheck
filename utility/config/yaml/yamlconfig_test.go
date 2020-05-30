package yaml_test

import (
	"gohealthcheck/utility/config/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

type yamlTestData struct {
	key string
	val string
}

var (
	fileName = "appsettings.yaml"

	yamlConfigTestSample = `token: "sampletoken"
bestnum: 42
`
	yamlConfigTestSampleData = []yamlTestData{
		{key: "token", val: "sampletoken"},
		{key: "bestnum", val: "42"},
	}
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
	err := ioutil.WriteFile(outFile, []byte(yamlConfigTestSample), os.ModePerm)
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestYamlConfig(c *C) {
	config, err := yaml.NewConfig(filepath.Join(s.dir, fileName))
	c.Assert(err, IsNil)

	dataIndex := 0
	for _, d := range yamlConfigTestSampleData {
		val := config.GetConfig(d.key)
		c.Assert(val, Equals, d.val)
		dataIndex++
	}
}
