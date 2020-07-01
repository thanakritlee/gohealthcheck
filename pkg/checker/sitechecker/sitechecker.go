package sitechecker

import (
	"errors"
	"gohealthcheck/pkg/checker"
	yamlConfig "gohealthcheck/pkg/config/yaml"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

var (
	// ErrorWrongCheckeeType error message.
	ErrorWrongCheckeeType = errors.New("wrong checkee type")
)

// SiteChecker implements the checker interface, and is
// use for checking the health status of a site.
type SiteChecker struct {
}

// NewSiteChecker return a new SiteChecker object.
func NewSiteChecker() *SiteChecker {
	return &SiteChecker{}
}

// Check check the health status of a site.
// Return true, if there're any http status code that is returned from the website.
// Return false, if it cannot reach the website (request timeout).
func (s *SiteChecker) Check(checkee checker.Checkee) (bool, error) {
	if checkee.Type != checker.TypeSite {
		return false, xerrors.Errorf("SiteChecker.Check: %w", ErrorWrongCheckeeType)
	}

	config := yamlConfig.Config{}
	timeOutConfig := config.GetConfig("timeout")
	if timeOutConfig == "" {
		timeOutConfig = "30"
	}
	timeOut, err := strconv.ParseInt(timeOutConfig, 10, 64)
	if err != nil {
		return false, xerrors.Errorf("SiteChecker.Check: %w", err)
	}
	httpClient := &http.Client{
		Timeout: time.Duration(timeOut * 1000000000),
	}

	_, err = httpClient.Get(checkee.URL)
	if err != nil {
		if (err.(*url.Error)).Timeout() {
			return false, nil
		}
		return false, xerrors.Errorf("SiteChecker.Check: %w", err)
	}

	return true, nil
}
