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
type IntervalSlice []*Interval

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
	return i[a].Less(i[b])
}

/* end sort.Interface methods */

func (i Interval) Stab(q uint16) bool {
	return i.Start <= q && q <= i.End
}

/*
*	Private
 */
type uniqueInterval struct {
	interval *Interval
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
type uniqueIntervalSlice []*uniqueInterval

func (a *uniqueInterval) Less(b *uniqueInterval) bool {
	// If they are equal, then we should also order by uniqueid
	// That will help ensure output maintains input order for
	// ranges that are the same
	aS := a.interval.Start
	aE := a.interval.End
	bS := b.interval.Start
	bE := b.interval.End
	return aS < bS ||
		(aS == bS &&
			(aE < bE || (aE == bE && a.id < b.id)))
	// This is dumb
}

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
	return i[a].Less(i[b])
}

/* end sort.Interface methods */

func (i uniqueInterval) Stab(q uint16) bool {
	return i.interval.Stab(q)
}
