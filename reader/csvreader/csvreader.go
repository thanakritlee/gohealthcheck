package csvreader

import (
	"encoding/csv"
	"errors"
	"gohealthcheck/reader"
	"os"

	"golang.org/x/xerrors"
)

var (
	// ErrorMissingCSVFile error message.
	ErrorMissingCSVFile = errors.New("missing CSV file")
)

// CSVReader implements the Reader interface.
type CSVReader struct {
	fileReader *csv.Reader
	file       *os.File
}

// NewCSVReader returns a new instance of the CSVReader which
// can be use to read the site data.
func NewCSVReader(csvFilePath string) (*CSVReader, error) {

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, xerrors.Errorf("NewCSVReader: %w", err)
	}

	return &CSVReader{fileReader: csv.NewReader(csvFile), file: csvFile}, nil
}

// ReadSite returns an Iterator for the sites.
func (c *CSVReader) ReadSite() (reader.SiteIterator, error) {
	return &siteIterator{csvReader: c.fileReader, file: c.file}, nil
}
