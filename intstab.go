/* integer stabbing. */
package intstab

import (
	"fmt"
	"log"
	"math"
	"sort"
)

/*
*	Public
 */
type Interval struct {
	Start uint16
	End   uint16
	Tag   interface{}
}
type IntervalSlice []Interval

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
	intervals IntervalSlice
	smaller   []IntervalSlice
}

func (i Interval) isValid() (ok bool, err error) {
	if i.Start > i.End {
		err = fmt.Errorf("Invalid interval: %d should be <= %d",
			i.Start, i.End)
	} else if i.Tag == nil {
		err = fmt.Errorf("Invalid Interval: Missing Tag")
	} else {
		ok = true
	}

	return
}

// Use a sweep line to build the tree in O(n)
func (ts *intStab) precompute() {
	// math.MaxUint16
	ts.precomputeSmaller()

	n := len(ts.intervals)
	event := make(map[uint16]IntervalSlice, n*2)
	list := make(IntervalSlice, n+1)

	// Sort by starting values
	sort.Sort(ts.intervals)
	for _, a := range ts.intervals {
		event[a.End] = append(event[a.End], a)
		event[a.Start] = append(event[a.Start], a)
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
	a, and every element in them is removed from I.
*/
func (ts *intStab) precomputeSmaller() {
	ts.smaller = make([]IntervalSlice, math.MaxUint16)
	sm := ts.smaller

	//log.Printf("Intervals %v", ts.intervals)
	for _, a := range ts.intervals {
		sm[a.Start] = append(sm[a.Start], a)
	}

	// sort Q elements by length
	for _, arr := range sm {
		if len(arr) > 1 {
			sort.Sort(arr)
		}

		if len(arr) > 0 {
			// remove the last one
			log.Print(arr)
			log.Printf("Val %v", arr[len(arr)-1])
			//arr[len(arr)-1] = nil
			arr = arr[:len(arr)-1]
			log.Print(arr)

			// Remove each remaining item from the intervals map
			for _, x := range arr {
				// Remove from ts.intervals
				log.Print("Should remove ", x)
			}
		}
	}
}

/* sort.Interface methods for []Interval */
func (i IntervalSlice) Len() int {
	return len(i)
}

func (i IntervalSlice) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

// for two intervals a and b, it holds a < b
// if la < lb or (la = lb ∧ ra ≤ rb).
// li = left (Start), ri = right (End)
func (i IntervalSlice) Less(a, b int) bool {
	return i[a].Start < i[b].Start ||
		(i[a].Start == i[b].Start && i[a].End <= i[b].End)
}

/* end sort.Interface methods */

func (ts *intStab) push(intervals IntervalSlice) (err error) {
	for _, val := range intervals {
		if _, err = val.isValid(); err != nil {
			return
		}
		ts.intervals = append(ts.intervals, val)
	}

	ts.precompute()
	return
}
