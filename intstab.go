/* integer stabbing. */
package intstab

import (
	//"log"
	"fmt"
)

/*
*	Public
 */
type Interval struct {
	Start int
	End   int
	Tag   interface{}
}

type IntervalStabber interface {
	Intersect(int) ([]Interval, bool)
}

func (ts *intStab) Intersect(query int) (results []Interval, ok bool) {
	results = make([]Interval, 10)

	return
}

func NewIntervalStabber(intervals ...Interval) (IntervalStabber, error) {
	ts := new(intStab)
	err := error(nil)

	err = ts.push(intervals...)
	if err != nil {
		return nil, err
	}

	return ts, err
}

/*
*	Private
 */
type intStab struct {
	intervals []Interval
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

func (ts *intStab) rebuildTree() {

}

func (ts *intStab) push(intervals ...Interval) (err error) {
	for _, val := range intervals {
		if _, err = val.isValid(); err != nil {
			return
		}
		ts.intervals = append(intervals, val)
	}

	ts.rebuildTree()

	return
}
