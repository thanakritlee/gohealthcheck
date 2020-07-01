package reader

// Iterator can be implemented by the Site objects so that it can be iterated.
type Iterator interface {
	Next() bool
	Error() error
	Close() error
	Total() (int, error)
}

// Site contains the URL of the site to check.
type Site struct {
	URL string
}

// SiteIterator can be implemented by Site objects so that it can be iterated.
type SiteIterator interface {
	Iterator
	Site() *Site
}

// Reader can be implemented by objects that can read sites from a source.
type Reader interface {
	ReadSite() (SiteIterator, error)
}
