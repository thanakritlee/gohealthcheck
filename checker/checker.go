package checker

// CheckeeType shows the types that can be check.
type CheckeeType uint8

const (
	// TypeSite check site's health.
	TypeSite CheckeeType = iota
)

// Checker can be implemented by objects that can check a checkee.
type Checker interface {
	Check(Checkee) (bool, error)
}

// Checkee states the type of checkee to check.
// I've implemented it this way, so that the method signature
// for Check won't have to be change in the future. Only the
// Checkee struct has to change if we want to add other things
// we want to check other than sites (URL).
type Checkee struct {
	Type CheckeeType
	URL  string
}
