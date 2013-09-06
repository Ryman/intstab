package intstab

type Stack interface {
	Push(*uniqueInterval)
	Pop() *uniqueInterval
	Len() int
}

func NewStack() Stack {
	s := new(sliceStack)
	s.data = make(uniqueIntervalSlice, 0, 250)
	return s
}

/*
*	Stack backed by slice
 */
type sliceStack struct {
	data uniqueIntervalSlice
}

func (s *sliceStack) Push(v *uniqueInterval) {
	s.data = append(s.data, v)
}

func (s *sliceStack) Pop() (v *uniqueInterval) {
	if l := len(s.data) - 1; l == -1 {
		return nil
	} else {
		// Pop(): https://code.google.com/p/go-wiki/wiki/SliceTricks
		v, s.data = s.data[l], s.data[:l]
		return
	}
}

func (s *sliceStack) Len() int {
	return len(s.data)
}
