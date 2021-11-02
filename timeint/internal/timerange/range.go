package timerange

import "time"

type Range struct {
	Start time.Time
	End   time.Time
}

// NewRange returns a new range representing the endpoints of a time-interval,
// panics if the right endpoint is before the left endpoint.
func NewRange(a, b time.Time) Range {
	if b.Before(a) {
		panic("invalid range")
	}
	return Range{Start: a, End: b}
}

func (r Range) empty() bool {
	return r.Start.Equal(r.End)
}

func (r Range) contains(t time.Time) bool {
	if r.empty() {
		return false
	}
	return r.Start.Before(t) && t.Before(r.End)
}

func (r Range) ClosedContains(t time.Time) bool {
	return r.Start.Equal(t) || r.contains(t) || t.Equal(r.End)
}

func (r Range) OpenContains(t time.Time) bool {
	return r.contains(t)
}

func (r Range) LeftClosedRightOpenContains(t time.Time) bool {
	return (!r.empty() && r.Start.Equal(t)) || r.contains(t)
}

func (r Range) LeftOpenRightClosedContains(t time.Time) bool {
	return r.contains(t) || (!r.empty() && t.Equal(r.End))
}
