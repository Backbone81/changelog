package changelog

import "time"

// Release is a struct holding all relevant information about a release and its changelog entries
type Release struct {
	Version string
	Date    time.Time
	Entries []Entry
}
