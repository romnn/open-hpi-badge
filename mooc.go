package openhpibadge

import (
	"time"
)

// MOOC ...
type MOOC struct {
	URL string
	Title string
	Participants Participants
	Language string
	Start time.Time
	End time.Time
}

// Participants ...
type Participants struct {
	Start, End, Current int
}

// Equal ...
func (m *MOOC) Equal(other *MOOC) (bool, bool) {
	equal := true
	equal = equal && m.URL == other.URL
	equal = equal && m.Title == other.Title
	equal = equal && m.Language == other.Language
	equal = equal && m.Start == other.Start
	equal = equal && m.End == other.End
	exact := equal && m.Participants == other.Participants
	return exact, equal
}