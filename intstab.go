/* integer stabbing. */
package intstab

import (
	_ "log"
	"math"
	"sort"
)

/*
*	Public
 */
type IntervalStabber interface {
	Intersect(uint16) ([]Interval, bool)
}

func (ts *intStab) Intersect(query uint16) (results []Interval, ok bool) {
	results = make([]Interval, 10)

	return
}

func NewIntervalStabber(intervals ...Interval) (IntervalStabber, error) {
	ts := new(intStab)
	err := error(nil)

	err = ts.push(intervals)
	if err != nil {
		return nil, err
	}

	return ts, err
}

/*
*	Private
 */
type intStab struct {
	intervals map[int]uniqueInterval
	smaller   []uniqueIntervalSlice
}

// Use a sweep line to build the tree in O(n)
func (ts *intStab) precompute() {
	// math.MaxUint16
	ts.precomputeSmaller()

	n := len(ts.intervals)
	event := make(map[uint16][]uniqueInterval, n*2)
	list := make(IntervalSlice, n+1)

	// Sort by starting values
	//sort.Sort(ts.intervals)
	for _, a := range ts.intervals {
		event[a.interval.End] = append(event[a.interval.End], a)
		event[a.interval.Start] = append(event[a.interval.Start], a)
	}

	for q := 0; q < n; q++ {
		if len(list) > 1 {

		}
	}
}

/*
	All the intervals with left endpoint l, except one such longest
	interval a, are stored in a list called Smaller(a) (see Figure 1).
	These lists are sorted by length in descending order, get a link to
	a, and every element in them is removed from I. */
func (ts *intStab) precomputeSmaller() {
	sm := make([]uniqueIntervalSlice, math.MaxUint16)

	//log.Printf("Intervals %v", ts.intervals)
	for _, a := range ts.intervals {
		sm[a.interval.Start] = append(sm[a.interval.Start], a)
	}

	// sort Q elements by length
	for _, arr := range sm {
		if len(arr) > 1 {
			sort.Sort(arr)
			/* Debug
			log.Print(arr)
			log.Printf("Val %v", arr[len(arr)-1]) */
			// remove the longest one (it should be the last one)
			arr = arr[:len(arr)-1]

			// Remove each remaining item from the intervals map
			for _, x := range arr {
				delete(ts.intervals, x.id)
			}
		} else if len(arr) == 1 {
			// There was only one, it's the longest so just remove and ignore
			arr = arr[:0]
		}
	}

	ts.smaller = sm
}

func (ts *intStab) push(intervals IntervalSlice) (err error) {
	ts.intervals = make(map[int]uniqueInterval, len(intervals))

	for id, val := range intervals {
		if _, err = val.isValid(); err != nil {
			return
		}

		ts.intervals[id] = uniqueInterval{val, id}
	}

	ts.precompute()
	return
}
