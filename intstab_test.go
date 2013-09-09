package intstab

import (
	"testing"
)

func dependOn(t *testing.T, deps ...func(*testing.T)) {
	for _, dependancy := range deps {
		dependancy(t)

		if t.Skipped() || t.Failed() {
			t.Skip("Dependency failed, Skipping.")
		}
	}
}

func checkExpected(t *testing.T, stab IntervalStabber, query uint16, expected ...interface{}) {
	actual, err := stab.Intersect(query)
	if err != nil {
		t.Fatalf("Error performing query: %v", err)
	}

	t.Log("Results:")
	for _, v := range actual {
		t.Log(v.Tag, " ")
	}
	t.Log("Expected:\n", expected)

	if l, x := len(actual), len(expected); l != x {
		t.Fatalf("Results were not expected length: Got %v but wanted %v", l, x)
	}

	for i, x := range expected {
		if actual[i].Tag != x {
			t.Errorf("Missing %s from results", x)
		}
	}
}

func setup(t *testing.T, intervals IntervalSlice) (stab IntervalStabber) {
	stab, err := NewIntervalStabber(intervals)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	return
}

func TestBasicQuery(t *testing.T) {
	intervals := IntervalSlice{{4, 15, "First"}, {34, 72, "Second"}}
	stab := setup(t, intervals)

	checkExpected(t, stab, 45, "Second")
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
		{0, 142, "Zero is cool"},
		{6000, 65535, "So is 65535"},
	}

	stab := setup(t, intervals)

	// Check we can query zero
	checkExpected(t, stab, 0, "Zero is cool")

	// Check we can query maxUint16
	// TODO: Fix this
	checkExpected(t, stab, 65535, "So is 65535")
}

func TestSmallOverlappingIntervals(t *testing.T) {
	intervals := IntervalSlice{
		{995, 995, "Copy"},
		{994, 995, "Another"},
		{995, 995, "Singular"},
		{989, 995, "Seventh"},
	}

	stab := setup(t, intervals)
	// TODO: Fix this
	checkExpected(t, stab, 995, "Seventh", "Another", "Copy", "Singular")
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

	stab := setup(t, intervals)

	checkExpected(t, stab, 42, "Fourth", "Seventh", "Third")
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
	// TestMultipleResultsQuery needs to work as expected
	dependOn(t, TestMultipleResultsQuery)

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

	stab := setup(t, intervals)

	allocs := testing.AllocsPerRun(1000, func() {
		results, err := stab.Intersect(42)
		if err != nil {
			t.Fatal("Error during alloc run: ", err)
		}

		if results == nil {
			t.Fatal("Got 'nil' results during alloc run")
		}
	})

	t.Log("Allocs per run (avg): ", allocs)
	if allocs > 2.1 {
		t.Fatal("Too many allocs, be sure to disable logging for real builds")
	}
}

func BenchmarkLargeQuery(b *testing.B) {
	intervals := IntervalSlice{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
		{34, 45, "Fourth"},
		{34, 40, "Fifth"},
		{34, 34, "Sixth"},
		{34, 45, "Seventh"},
	}
	// Make a result for every q
	for i := 0; i < 65535; i++ {
		intervals = append(intervals, &Interval{uint16(i), uint16(i), i})
	}

	stab, err := NewIntervalStabber(intervals)
	if err != nil {
		b.Fatalf("Unable to setup benchmark: ", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results, err := stab.Intersect(42)
		if err != nil || len(results) != 4 {
			b.Fatal("Tests failed during benchmark")
		}
	}
}

func BenchmarkQuery(b *testing.B) {
	intervals := IntervalSlice{
		{4, 15, "First"},
		{50, 72, "Second"},
		{34, 90, "Third"},
		{34, 45, "Fourth"},
		{34, 40, "Fifth"},
		{34, 34, "Sixth"},
		{34, 45, "Seventh"},
	}

	stab, err := NewIntervalStabber(intervals)
	if err != nil {
		b.Fatalf("Unable to setup benchmark: ", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results, err := stab.Intersect(42)
		if err != nil || len(results) != 3 {
			b.Fatal("Tests failed during benchmark")
		}
	}
}
