package csvreader

import (
	"encoding/csv"
	"errors"
	"gohealthcheck/reader"
	"os"

	"golang.org/x/xerrors"
)

type siteIterator struct {
	file      *os.File
	csvReader *csv.Reader
	site      *reader.Site
	lastError error
}

func (i *siteIterator) Next() bool {
	row, err := i.csvReader.Read()
	if err != nil {
		i.lastError = xerrors.Errorf("siteIterator.Next: %w", err)
		return false
	}

	if len(row) > 0 {
		i.site = &reader.Site{URL: row[0]}
	} else {
		i.site = &reader.Site{URL: ""}
	}

	return true
}

func (i *siteIterator) Error() error {
	return i.lastError
}

func (i *siteIterator) Close() error {
	return i.file.Close()
}

func (i *siteIterator) Site() *reader.Site {
	site := new(reader.Site)
	*site = *i.site
	return site
}

func (i *siteIterator) Total() (int, error) {
	return 0, xerrors.Errorf("siteIterator.Total: %w", errors.New("CSV Reader isn't able to get total row"))
}
