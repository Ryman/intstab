package intstab

import (
	"testing"
)

func TestBasicQuery(t *testing.T) {
	// Test intervals
	intervals := IntervalSlice{{4, 15, "First"}, {34, 72, "Second"}}

	ts, err := NewIntervalStabber(intervals)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	results, err := ts.Intersect(45)
	if err != nil {
		t.Fatalf("Incorrect")
	}

	if len(results) != 1 {
		t.Fatal("Wrong number of results for Intersect")
	}

	if results[0].Tag != "Second" {
		t.Error("Wrong result from Intersect")
	}

	if t.Failed() {
		t.Log("Results were: ", results)
	}
}

func TestNoResultsQuery(t *testing.T) {
	intervals := IntervalSlice{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
		{34, 45, "Fourth"},
		{34, 40, "Fifth"},
		{34, 34, "Sixth"},
		{34, 45, "Seventh"},
	}

	ts, err := NewIntervalStabber(intervals)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	// Check we can query zero
	results, err := ts.Intersect(0)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if len(results) > 0 {
		t.Error("Unexpected results: ", results)
	}

	// Check we can query maxUint16
	results, err = ts.Intersect(65535)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if len(results) > 0 {
		t.Error("Unexpected results: ", results)
	}
}

func TestMultipleResultsQuery(t *testing.T) {
	// Test intervals
	intervals := IntervalSlice{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
		{34, 45, "Fourth"},
		{34, 40, "Fifth"},
		{34, 34, "Sixth"},
		{34, 45, "Seventh"},
	}

	ts, err := NewIntervalStabber(intervals)
	if err != nil {
		t.Fatal(err)
	}

	results, err := ts.Intersect(42)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if len(results) != 3 {
		t.Fatal("Wrong number of results for Intersect: %v", results)
	}

	if results[0].Tag != "Fourth" {
		// Ensure the resultant ordering is ordered by leftmost interval.Start
		t.Error("Wrong result from Intersect")
	}

	if results[1].Tag != "Seventh" {
		// Ensure we get multiple different results for the same range
		// We also need to ensure the ordering is the same as it went in
		t.Error("Missing an overlapping range for Intersect")
	}

	if results[2].Tag != "Third" {
		t.Error("Wrong result from Intersect")
	}

	if t.Failed() {
		t.Log("Results were: ", results)
	}
}

func TestIntervalBadRange(t *testing.T) {
	// Test intervals
	intervals := IntervalSlice{{4, 15, "First"}, {340, 72, "Second"}}

	_, err := NewIntervalStabber(intervals)
	if err == nil {
		t.Fatalf("Should not have accepted invalid interval")
	}
}

func TestIntervalBadTag(t *testing.T) {
	// Test intervals
	intervals := IntervalSlice{{4, 15, "First"}, {34, 72, nil}}

	_, err := NewIntervalStabber(intervals)
	if err == nil {
		t.Fatalf("Should not have accepted nil Tag for interval")
	}
}

func TestIntervalSorting(t *testing.T) {
	// Test by Start values

	// Test equal Start values but different Ends
}

func TestOptimalTime(t *testing.T) {
	// Make a tiny dataset
	// Gather a time x for y results

	// Make a massive dataset
	// Show that a query with y results is similar to time x

	// Show that y < 2y < 4y results
}

func TestGarbageCreation(t *testing.T) {
	// Test intervals
	intervals := IntervalSlice{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
		{34, 45, "Fourth"},
		{34, 40, "Fifth"},
		{34, 34, "Sixth"},
		{34, 45, "Seventh"},
	}

	ts, err := NewIntervalStabber(intervals)
	if err != nil {
		t.Fatal(err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		results, err := ts.Intersect(42)
		if err != nil {
			t.Fatal("Error during alloc run: ", err)
		}

		if results == nil {
			t.Fatal("Got 'nil' results during alloc run")
		}
	})

	t.Log("Allocs per run (avg): ", allocs)
	if allocs > 2.5 {
		t.Fatal("Too many allocs, be sure to disable logging for real builds")
	}
}
