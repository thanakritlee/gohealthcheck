package console

import (
	"fmt"
	"gohealthcheck/pkg/reporter"
)

// Reporter implements the reporter.Reporter interface.
type Reporter struct{}

// NewConsoleReporter return a new console-Reporter instance.
func NewConsoleReporter() *Reporter {
	return &Reporter{}
}

// Report logs the statistics to the console.
func (r *Reporter) Report(stats reporter.Statistics) error {
	fmt.Printf("Checked websites: [%d]\n", stats.TotalWebSites)
	fmt.Printf("Successful websites: [%d]\n", stats.Success)
	fmt.Printf("Failure websites: [%d]\n", stats.Failure)
	fmt.Printf("Total time to finished checking websites: [%d ns]\n", stats.TotalTime)

	return nil
}
