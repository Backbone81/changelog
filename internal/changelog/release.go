package changelog

import "time"

type Release struct {
	Version string
	Date    time.Time
	Entries []Entry
}
