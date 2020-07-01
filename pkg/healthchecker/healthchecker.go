package healthchecker

import (
	"gohealthcheck/pkg/checker"
	"gohealthcheck/pkg/checker/sitechecker"
	"gohealthcheck/pkg/reader"
	"gohealthcheck/pkg/reader/csvreader"
	"gohealthcheck/pkg/reporter"
	"gohealthcheck/pkg/reporter/console"
	httpReporter "gohealthcheck/pkg/reporter/http"
	"net/http"
	"sync"
	"time"

	"golang.org/x/xerrors"
)

// InputType describes the type of the data input method.
type InputType uint8

// ReportType describes the type of the report.
type ReportType uint8

const (
	// TypeHTTP http report type.
	TypeHTTP ReportType = iota
	// TypeConsole console report type.
	TypeConsole

	// TypeCSV csv data input type.
	TypeCSV InputType = iota
)

// Input describe the data input method.
type Input struct {
	InputType   InputType
	CSVFileName string
}

// Report describes how the results are reported.
type Report struct {
	ReportType ReportType
	URL        string
	Headers    http.Header
}

// HCInterface  can be implemented by objects that can check healths of something.
type HCInterface interface {
	CheckHealth() error
}

// HC implements the HealthCheckerInf interface.
type HC struct {
	input        Input
	report       []Report
	siteIterator reader.SiteIterator
	stats        reporter.Statistics
}

// NewHealthChecker returns a new HC instance.
func NewHealthChecker(input *Input, report *[]Report) (*HC, error) {

	hc := &HC{}

	hcInput := new(Input)
	*hcInput = *input

	hcReport := new([]Report)
	*hcReport = *report

	hc.input = *hcInput
	hc.report = *hcReport

	hc.importData()

	return hc, nil

}

// CheckHealth check the health of the sites.
func (h *HC) CheckHealth() error {

	startTime := time.Now()

	totalWebsites := 0
	numberOfSuccess := 0
	numberOfFailure := 0

	var checkWG sync.WaitGroup

	var resultWG sync.WaitGroup
	resultWG.Add(1)
	resultChan := make(chan bool)

	siteChecker := sitechecker.NewSiteChecker()

	for (h.siteIterator).Next() {
		checkWG.Add(1)
		site := h.siteIterator.Site()

		totalWebsites++

		go func(url string) {
			defer checkWG.Done()

			checkee := checker.Checkee{
				Type: checker.TypeSite,
				URL:  url,
			}

			// Timeout or cannot reach the site in any way
			// is considered a failure.
			success, _ := siteChecker.Check(checkee)

			resultChan <- success

		}(site.URL)

	}

	go func() {
		defer resultWG.Done()

		for {
			success, ok := <-resultChan
			if !ok {
				return
			}

			if success {
				numberOfSuccess++
			} else {
				numberOfFailure++
			}
		}
	}()

	checkWG.Wait()
	close(resultChan)
	resultWG.Wait()

	stopTime := time.Now()
	totalTime := stopTime.Sub(startTime).Nanoseconds()

	h.stats = reporter.Statistics{
		TotalWebSites: totalWebsites,
		Success:       numberOfSuccess,
		Failure:       numberOfFailure,
		TotalTime:     totalTime,
	}

	return nil

}

// ReportResult reports the health check result using the reporter object.
// The reporter object implementation depends on the report type that were
// specified when instantiating the HealthChecker object.
func (h *HC) ReportResult() error {
	var reporter reporter.Reporter
	var err error

	for _, reportType := range h.report {
		switch reportType.ReportType {
		case TypeHTTP:
			reporter = httpReporter.NewHTTPReporter(reportType.URL, reportType.Headers)
			err = reporter.Report(h.stats)
		case TypeConsole:
			reporter = console.NewConsoleReporter()
			err = reporter.Report(h.stats)
		}

		if err != nil {
			return xerrors.Errorf("HC.ReportResult: %w", err)
		}
	}

	return nil
}

func (h *HC) importData() error {
	switch h.input.InputType {
	case TypeCSV:
		csvReader, err := csvreader.NewCSVReader(h.input.CSVFileName)
		if err != nil {
			return xerrors.Errorf("HC.importData: %w", err)
		}

		h.siteIterator, err = csvReader.ReadSite()
		if err != nil {
			return xerrors.Errorf("HC.importData: %w", err)
		}
	}

	return nil
}
