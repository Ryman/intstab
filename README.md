intstab
=======

A Golang Interval stabbing implementation for small integer ranges

As described in "Interval Stabbing Problems in Small Integer Ranges" Jens M. Schmidt, 2010


Assuming a set of intervals, I, it will find the answer to which intervals cover query q
in o(1+k) time where k is the result set size. The output will be ordered.

Example usage:

    intervals := IntervalSlice{
		  {4, 15, "First"},
		  {50, 72, "Second"},
  		{34, 90, "Third"},
  		{34, 45, "Fourth"},
  		{34, 40, "Fifth"},
  		{34, 34, "Sixth"},
  		{34, 45, "Seventh"},
  	}

    // Initialise
	  ts, _ := NewIntervalStabber(intervals)

    // Query for intervals intersecting 42
	  results, _ := ts.Intersect(42)
	  
	  // Results should be:
	  // [0].Tag = "Fourth"
	  // [1].Tag = "Seventh"
	  // [2].Tag = "Third"

