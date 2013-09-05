package intstab

import (
	"fmt"
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

// for two intervals a and b, it holds a < b
// if la < lb or (la = lb ∧ ra ≤ rb).
// li = left (Start), ri = right (End)
func (a *Interval) Less(b *Interval) bool {
	return a.Start < b.Start ||
		(a.Start == b.Start && a.End <= b.End)
}

/* sort.Interface methods for []Interval */
func (i IntervalSlice) Len() int {
	return len(i)
}

func (i IntervalSlice) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i IntervalSlice) Less(a, b int) bool {
	return i[a].Less(&i[b])
}

/* end sort.Interface methods */

/*
*	Private
 */
type uniqueInterval struct {
	interval Interval
	id       int
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

// TODO: Probably a more idomatic way to do this
type uniqueIntervalSlice []uniqueInterval

/* sort.Interface methods for []Interval */
func (i uniqueIntervalSlice) Len() int {
	return len(i)
}

func (i uniqueIntervalSlice) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

// for two intervals a and b, it holds a < b
// if la < lb or (la = lb ∧ ra ≤ rb).
// li = left (Start), ri = right (End)
func (i uniqueIntervalSlice) Less(a, b int) bool {
	return i[a].interval.Less(&i[b].interval)
}

/* end sort.Interface methods */
