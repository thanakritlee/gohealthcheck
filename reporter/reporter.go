package reporter

// Statistics data to report.
type Statistics struct {
	TotalWebSites int   `json:"total_websites"`
	Success       int   `json:"success"`
	Failure       int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

// Reporter can be implemented by objects that can report the statistic.
type Reporter interface {
	Report(Statistics) error
}
