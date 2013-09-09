/* integer stabbing. */
package intstab

import (
	"container/list"
	"fmt"
	"math"
	"sort"
)

/*
*	Public
 */
type IntervalStabber interface {
	Intersect(uint16) ([]*Interval, error)
}

func (ts *intStab) Intersect(query uint16) (results []*Interval, err error) {
	return ts.stab(query)
}

func NewIntervalStabber(intervals IntervalSlice) (ts IntervalStabber, err error) {
	ts = new(intStab)
	err = ts.(*intStab).init(intervals)

	if err != nil {
		ts = nil
	}

	return
}

/*
*	Private
 */

/*
* +1 to always allocate including zero (allows queries 0-65535 inclusive)
* Anything higher and the compiler stops it =)
 */
const maxN = math.MaxUint16 + 1

type intStab struct {
	root      *uniqueInterval
	intervals map[int]*uniqueInterval
	starts    uniqueIntervalSlice
}

func (ts *intStab) init(intervals IntervalSlice) (err error) {
	n := len(intervals)
	ts.intervals = make(map[int]*uniqueInterval, n)

	for id, val := range intervals {
		if _, err = val.isValid(); err != nil {
			return
		}

		ts.intervals[id] = &uniqueInterval{interval: val, id: id}
	}

	// Root is id(-1)
	ts.root = &uniqueInterval{interval: &Interval{}, id: -1}

	ts.precompute()
	return
}

func (ts *intStab) precompute() {
	ts.precomputeSmaller()
	events := ts.buildEvents()
	starts := make(uniqueIntervalSlice, maxN)
	children := make(map[int]*list.List, len(ts.intervals)) // include space for root's children

	l := list.New()
	lmap := make(map[int]*list.Element)

	// Use a sweep line to build the tree in O(n)
	for q := 0; q < maxN; q++ {
		if v := l.Back(); v != nil {
			starts[q] = v.Value.(*uniqueInterval)
		}

		// For all events in reverse order
		for i := len(events[q]) - 1; i >= 0; i-- {
			a := events[q][i]

			// Rightmost left
			if int(a.interval.Start) == q &&
				(starts[q] == nil ||
					starts[q].id != a.id) { // Catch single ranges e.g. 45-45
				starts[q] = a
				// Save link to position in L
				lmap[a.id] = l.PushBack(a)
			} else { // a.interval.End == q
				var parent *uniqueInterval
				e := lmap[a.id]

				if pred := e.Prev(); pred != nil {
					// a has predecessor b
					parent = pred.Value.(*uniqueInterval)
				} else {
					// root as Parent
					parent = ts.root
				}

				// Calculate siblings
				if kids := children[parent.id]; kids == nil {
					children[parent.id] = list.New()
				} else if leftChild := kids.Front(); leftChild != nil {

					a.leftSibling = leftChild.Value.(*uniqueInterval)
				}

				a.parent = parent
				children[parent.id].PushBack(a)

				l.Remove(e)
			}
		}
	}

	ts.starts = starts
}

func (ts *intStab) buildEvents() []uniqueIntervalSlice {
	events := make([]uniqueIntervalSlice, maxN)

	// Get a sorted list from the intervals map
	sorted := make(uniqueIntervalSlice, 0, len(ts.intervals))
	for _, v := range ts.intervals {
		sorted = append(sorted, v)
	}

	// Sort by starting values
	sort.Sort(sorted)

	// Build the events
	for _, a := range sorted {
		events[a.interval.End] = append(events[a.interval.End], a)
		events[a.interval.Start] = append(events[a.interval.Start], a)
	}

	return events
}

/*
	All the intervals with left endpoint l, except one such longest
	interval a, are stored in a list called Smaller(a) (see Figure 1).
	These lists are sorted by length in descending order, get a link to
	a, and every element in them is removed from I. */
func (ts *intStab) precomputeSmaller() {
	sm := make([]uniqueIntervalSlice, maxN)

	for _, a := range ts.intervals {
		sm[a.interval.Start] = append(sm[a.interval.Start], a)
	}

	// sort Q elements by length
	for i, arr := range sm {
		if l := len(arr); l > 1 {
			// More than one interval with same start, link them together
			// and keep the longest. Remove all the rest from interval list
			sort.Sort(arr)

			// remove the longest one (it should be the last one)
			// set it's nextSmaller before removing
			arr[l-1].nextSmaller = arr[l-2]
			arr = arr[:l-1]

			// Remove each remaining item from the intervals map
			for i, x := range arr {
				if i == 0 {
					x.nextSmaller = nil
				} else {
					x.nextSmaller = arr[i-1]
				}
				delete(ts.intervals, x.id)
			}
		}

		sm[i] = nil
	}
}

func (ts *intStab) stab(q uint16) (results []*Interval, err error) {
	if ts.starts == nil {
		return nil, fmt.Errorf("Need to initialise before querying")
	} else if len(ts.starts) <= int(q) {
		return nil, fmt.Errorf("Query is out of range")
	}

	stack := NewStack()

	for x := ts.starts[q]; x != nil && x.id != ts.root.id; x = x.parent {
		ts.traverse(x, stack, q)
	}

	finalLen := stack.Len()
	results = make(IntervalSlice, finalLen)

	for i := 0; i < finalLen; i++ {
		results[i] = stack.Pop().(*uniqueInterval).interval
	}

	return
}

func (ts *intStab) traverse(v *uniqueInterval, stack Stack, q uint16) {
	stack.Push(v)

	// Iterate over ts.smaller from largest to smallest
	for w := v.nextSmaller; w != nil; w = w.nextSmaller {
		if w.interval.End < q {
			break
		}
		stack.Push(w)
	}

	/* Calculate rightmost path:
	* The rightmost path R(v) of a node v ∈ V (S) is empty if v has no left
	* sibling or its left sibling w is not stabbed. Otherwise, R(v) is the
	* path from w to the rightmost stabbed node in the subtree of w in S. */
	for next := v.leftSibling; next != nil; next = next.leftSibling {
		// Check left siblings until they don't stab
		if next.Stab(q) {
			ts.traverse(next, stack, q)
		} else {
			break
		}
	}
}
