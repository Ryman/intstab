package intstab

type Stack interface {
	Push(interface{})
	Pop() interface{}
	Len() int
}

func NewStack() Stack {
	return new(PooledSliceStack)
}

/*
*	Stack backed by slices from a leaky pool
*   http://golang.org/doc/effective_go.html#leaky_buffer
* 	(The leaky buffer gave 33% perf improvement for test usage)
 */
type poolElement interface{}
type poolElementSlice []poolElement
type PooledSliceStack struct {
	data poolElementSlice
}

// default to 100 buffers, could probably reduce this
var freeList = make(chan poolElementSlice, 100)

func (s *PooledSliceStack) Push(v interface{}) {
	// Get a buffer from the pool if available
	if s.data == nil {
		select {
		case s.data = <-freeList:
			// Got a buffer, we're ok
		default:
			// No buffer free, alloc one of decent size
			s.data = make(poolElementSlice, 0, 250)
		}
	}
	s.data = append(s.data, v)
}

func (s *PooledSliceStack) Pop() (v interface{}) {
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
		return
	}
}

func (s *PooledSliceStack) Len() int {
	return len(s.data)
}
