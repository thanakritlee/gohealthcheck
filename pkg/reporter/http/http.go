package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"gohealthcheck/pkg/reporter"
	"net/http"

	"golang.org/x/xerrors"
)

// Reporter implements the reporter.Reporter interface.
type Reporter struct {
	url     string
	headers http.Header
}

// NewHTTPReporter returns a new http-Reporter instance.
func NewHTTPReporter(url string, headers http.Header) *Reporter {
	return &Reporter{
		url:     url,
		headers: headers,
	}
}

// Report sends the statistics to the HTTP endpoint.
func (r *Reporter) Report(stats reporter.Statistics) error {

	reqBody, err := json.Marshal(stats)
	if err != nil {
		return xerrors.Errorf("Reporter.Report: %w", err)
	}

	request, err := http.NewRequest("POST", r.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return xerrors.Errorf("Reporter.Report: %w", err)
	}
	request.Header = r.headers

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return xerrors.Errorf("Reporter.Report: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return xerrors.Errorf("Reporter.Report: %w", errors.New("failed to send report to endpoint"))
	}

	return nil
}
