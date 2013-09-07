/* integer stabbing. */
package intstab

import (
	"container/list"
	"fmt"
	"log"
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

func NewIntervalStabber(intervals IntervalSlice) (IntervalStabber, error) {
	ts := new(intStab)
	err := error(nil)

	err = ts.init(intervals)
	if err != nil {
		return nil, err
	}

	return ts, err
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
	children  map[int]*list.List
	parents   map[int]*uniqueInterval

	smaller []uniqueIntervalSlice
	starts  uniqueIntervalSlice
}

func (ts *intStab) init(intervals IntervalSlice) (err error) {
	n := len(intervals)
	ts.intervals = make(map[int]*uniqueInterval, n)
	ts.parents = make(map[int]*uniqueInterval, n)
	ts.children = make(map[int]*list.List, n+1) // include space for root's children

	for id, val := range intervals {
		if _, err = val.isValid(); err != nil {
			return
		}

		ts.intervals[id] = &uniqueInterval{val, id}
		ts.children[id] = list.New()
	}

	// Root is id(-1)
	ts.root = &uniqueInterval{&Interval{}, -1}
	ts.children[ts.root.id] = list.New()

	ts.precompute()
	return
}

func (ts *intStab) precompute() {
	ts.precomputeSmaller()
	events := ts.buildEvents()
	starts := make(uniqueIntervalSlice, maxN)

	l := list.New()
	lmap := make(map[int]*list.Element)

	// Use a sweep line to build the tree in O(n)
	for q := uint16(0); q < maxN-1; q++ {
		if v := l.Back(); v != nil {
			starts[q] = v.Value.(*uniqueInterval)
		}

		// For all events in reverse order
		for i := len(events[q]) - 1; i >= 0; i-- {
			a := events[q][i]

			// Rightmost left
			if a.interval.Start == q {
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

				ts.parents[a.id] = parent
				ts.children[parent.id].PushBack(a)

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

	//log.Printf("Intervals %v", ts.intervals)
	for _, a := range ts.intervals {
		sm[a.interval.Start] = append(sm[a.interval.Start], a)
	}

	// sort Q elements by length
	for _, arr := range sm {
		if len(arr) > 1 {
			sort.Sort(arr)
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

func (ts *intStab) stab(q uint16) (results []*Interval, err error) {
	if ts.starts == nil {
		err = fmt.Errorf("Need to initialise before querying")
	} else if len(ts.starts) <= int(q) {
		err = fmt.Errorf("Query is out of range")
	}

	if err != nil {
		return
	}

	x := ts.starts[q]
	if x == nil || x.interval == nil {
		return
	}

	stack := NewStack()

	for {
		ts.traverse(x, stack, q)
		x = ts.parents[x.id]

		if x == nil || x.id == ts.root.id {
			break
		}
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

	start := v.interval.Start
	// Iterate over ts.smaller from largest to smallest
	for i := len(ts.smaller[start]) - 1; i >= 0; i-- {
		w := ts.smaller[start][i]

		// Skip any interval larger than the current stabbed
		// it should already be in the stabbed list
		if w.interval.End >= v.interval.End {
			continue
		}

		// Check for intersection, we know q intersects something longer
		// If it doesn't, then no need to check smaller ones
		if w.interval.End < q {
			break
		}
		stack.Push(w) // Else, Add to stack
	}

	/* Calculate rightmost path:
	* The rightmost path R(v) of a node v ∈ V (S) is empty if v has no left
	* sibling or its left sibling w is not stabbed. Otherwise, R(v) is the
	* path from w to the rightmost stabbed node in the subtree of w in S. */
	parent := ts.parents[v.id]
	if parent == nil {
		parent = ts.root
	}
	children := ts.children[parent.id]

	// Find thyself
	// TODO: Make this O(1), haven't seen it get too high yet though
	if len := children.Len(); len > 2 {
		log.Print(children.Len())
	}

	var e *list.Element
	for e = children.Back(); e != nil; e = e.Prev() {
		// do something with e.Value
		if e.Value.(*uniqueInterval).id == v.id {
			break
		}
	}

	if e == nil {
		return
	}

	// Check left siblings until they don't stab
	for e = e.Prev(); e != nil; e = e.Prev() {
		if w := e.Value.(*uniqueInterval); w.Stab(q) {
			ts.traverse(w, stack, q)
		} else {
			break
		}
	}
}
