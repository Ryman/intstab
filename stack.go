package intstab

type UniqueIntervalStack interface {
	Push(*uniqueInterval)
	Pop() *uniqueInterval
	Len() int
}

func NewStack() UniqueIntervalStack {
	return new(sliceStack)
}

/*
*	Stack backed by slices from a leaky pool
*   http://golang.org/doc/effective_go.html#leaky_buffer
 */
type sliceStack struct {
	data uniqueIntervalSlice
}

// default to 100 buffers, could probably reduce this
var freeList = make(chan uniqueIntervalSlice, 100)

func (s *sliceStack) Push(v *uniqueInterval) {
	// Get a buffer from the pool if available
	if s.data == nil {
		select {
		case s.data = <-freeList:
			// Got a buffer, we're ok
		default:
			// No buffer free, alloc one of decent size
			s.data = make(uniqueIntervalSlice, 0, 250)
		}
	}
	s.data = append(s.data, v)
}

func (s *sliceStack) Pop() (v *uniqueInterval) {
	if l := len(s.data) - 1; l == -1 {
		return nil
	} else {
		// Pop(): https://code.google.com/p/go-wiki/wiki/SliceTricks
		v, s.data = s.data[l], s.data[:l]

		// Check if we're now empty, if so return resources to our bufferpool
		// Note: l is already -1 from above
		if l == 0 {
			select {
			case freeList <- s.data:
				// successfully put back into the pool
			default:
				// couldn't put it back in pool (full channel), just let it gc
			}
			// Remove our reference
			s.data = nil
		}
		s = nil
		return
	}
}

func (s *sliceStack) Len() int {
	return len(s.data)
}
