package timerange

import (
	"fmt"
	"testing"
	"time"
)

var (
	t0 = time.Unix(0, 0).In(time.UTC)
	t1 = t0.Add(time.Second)

	t0t0Range = Range{t0, t0}
	t0t1Range = Range{t0, t1}
)

func ExampleRange() {
	fmt.Println(Range{})
	// Output:
	// {0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}
}

func ExampleNewRange() {
	cases := []struct {
		a, b time.Time
	}{
		{t0, t0},
		{t0, t1},
		{t1, t1},
	}
	for _, c := range cases {
		r := NewRange(c.a, c.b)
		fmt.Println(r)
	}
	// Output:
	// {1970-01-01 00:00:00 +0000 UTC 1970-01-01 00:00:00 +0000 UTC}
	// {1970-01-01 00:00:00 +0000 UTC 1970-01-01 00:00:01 +0000 UTC}
	// {1970-01-01 00:00:01 +0000 UTC 1970-01-01 00:00:01 +0000 UTC}
}

func TestNewRange_WithRightBeforeLeft_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()
	NewRange(t1, t0)
}

type (
	containsFunc       func(time.Time) bool
	containsFuncChoice int
)

const (
	closed containsFuncChoice = iota
	open
	leftClosedRightOpen
	leftOpenRightClosed
)

func (c containsFuncChoice) asFunc(r Range, t time.Time) containsFunc {
	switch c {
	case closed:
		return func(t time.Time) bool {
			return r.ClosedContains(t)
		}
	case open:
		return func(t time.Time) bool {
			return r.OpenContains(t)
		}
	case leftClosedRightOpen:
		return func(t time.Time) bool {
			return r.LeftClosedRightOpenContains(t)
		}
	case leftOpenRightClosed:
		return func(t time.Time) bool {
			return r.LeftOpenRightClosedContains(t)
		}
	default:
		panic("invalid choice")
	}
}

func TestRangeContains(t *testing.T) {
	cases := []struct {
		r            Range
		t            time.Time
		containsFunc containsFuncChoice
		want         bool
	}{
		{t0t0Range, t0.Add(-1 * time.Nanosecond), closed, false},
		{t0t0Range, t0.Add(+0 * time.Nanosecond), closed, true},
		{t0t0Range, t0.Add(+1 * time.Nanosecond), closed, false},

		{t0t0Range, t0.Add(-1 * time.Nanosecond), open, false},
		{t0t0Range, t0.Add(+0 * time.Nanosecond), open, false},
		{t0t0Range, t0.Add(+1 * time.Nanosecond), open, false},

		{t0t0Range, t0.Add(-1 * time.Nanosecond), leftClosedRightOpen, false},
		{t0t0Range, t0.Add(+0 * time.Nanosecond), leftClosedRightOpen, false},
		{t0t0Range, t0.Add(+1 * time.Nanosecond), leftClosedRightOpen, false},

		{t0t0Range, t0.Add(-1 * time.Nanosecond), leftOpenRightClosed, false},
		{t0t0Range, t0.Add(+0 * time.Nanosecond), leftOpenRightClosed, false},
		{t0t0Range, t0.Add(+1 * time.Nanosecond), leftOpenRightClosed, false},

		{t0t1Range, t0.Add(-1 * time.Nanosecond), closed, false},
		{t0t1Range, t0.Add(+0 * time.Nanosecond), closed, true},
		{t0t1Range, t0.Add(+1 * time.Nanosecond), closed, true},
		{t0t1Range, t1.Add(-1 * time.Nanosecond), closed, true},
		{t0t1Range, t1.Add(+0 * time.Nanosecond), closed, true},
		{t0t1Range, t1.Add(+1 * time.Nanosecond), closed, false},

		{t0t1Range, t0.Add(-1 * time.Nanosecond), open, false},
		{t0t1Range, t0.Add(+0 * time.Nanosecond), open, false},
		{t0t1Range, t0.Add(+1 * time.Nanosecond), open, true},
		{t0t1Range, t1.Add(-1 * time.Nanosecond), open, true},
		{t0t1Range, t1.Add(+0 * time.Nanosecond), open, false},
		{t0t1Range, t1.Add(+1 * time.Nanosecond), open, false},

		{t0t1Range, t0.Add(-1 * time.Nanosecond), leftClosedRightOpen, false},
		{t0t1Range, t0.Add(+0 * time.Nanosecond), leftClosedRightOpen, true},
		{t0t1Range, t0.Add(+1 * time.Nanosecond), leftClosedRightOpen, true},
		{t0t1Range, t1.Add(-1 * time.Nanosecond), leftClosedRightOpen, true},
		{t0t1Range, t1.Add(+0 * time.Nanosecond), leftClosedRightOpen, false},
		{t0t1Range, t1.Add(+1 * time.Nanosecond), leftClosedRightOpen, false},

		{t0t1Range, t0.Add(-1 * time.Nanosecond), leftOpenRightClosed, false},
		{t0t1Range, t0.Add(+0 * time.Nanosecond), leftOpenRightClosed, false},
		{t0t1Range, t0.Add(+1 * time.Nanosecond), leftOpenRightClosed, true},
		{t0t1Range, t1.Add(-1 * time.Nanosecond), leftOpenRightClosed, true},
		{t0t1Range, t1.Add(+0 * time.Nanosecond), leftOpenRightClosed, true},
		{t0t1Range, t1.Add(+1 * time.Nanosecond), leftOpenRightClosed, false},
	}
	for _, c := range cases {
		f := c.containsFunc.asFunc(c.r, c.t)
		if got := f(c.t); got != c.want {
			t.Errorf("got %v but wanted %v: %+v", got, c.want, c)
		}
	}
}
