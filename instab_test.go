package intstab

import (
	"testing"
)

func TestIntStabInit(t *testing.T) {
	// Test intervals
	intervals := []Interval{{4, 15, "First"}, {34, 72, "Second"}}

	ts, err := NewIntervalStabber(intervals...)
	if err != nil {
		t.Fatal(err)
	}

	if results, ok := ts.Intersect(45); ok {
		t.Fatalf("Incorrect")
	} else if len(results) != 1 {
		t.Fatalf("Wrong number of results for Intersect")
	} else if results[0].Tag != "Second" {
		t.Fatalf("Wrong result from Intersect")
	}
}

func TestIntervalMultipleResults(t *testing.T) {
	// Test intervals
	intervals := []Interval{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
	}

	ts, err := NewIntervalStabber(intervals...)
	if err != nil {
		t.Fatal(err)
	}

	if results, ok := ts.Intersect(60); ok {
		t.Fatalf("Incorrect")
	} else if len(results) != 2 {
		t.Fatalf("Wrong number of results for Intersect")
	} else if results[0].Tag != "Third" {
		// Ensure the resultant ordering is ordered by leftmost interval.Start
		t.Fatalf("Wrong result from Intersect")
	} else if results[1].Tag != "Second" {
		t.Fatalf("Wrong result from Intersect")
	}
}

func TestIntervalBadRange(t *testing.T) {
	// Test intervals
	intervals := []Interval{{4, 15, "First"}, {340, 72, "Second"}}

	_, err := NewIntervalStabber(intervals...)
	if err == nil {
		t.Fatalf("Should not have accepted invalid interval")
	}
}

func TestIntervalBadTag(t *testing.T) {
	// Test intervals
	intervals := []Interval{{4, 15, "First"}, {34, 72, nil}}

	_, err := NewIntervalStabber(intervals...)
	if err == nil {
		t.Fatalf("Should not have accepted nil Tag for interval")
	}
}
