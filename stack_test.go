package intstab

import "testing"

func checkStackPopOrder(t *testing.T, s Stack, expected ...interface{}) {
	for _, val := range expected {
		if got := s.Pop(); got != val {
			t.Errorf("Unexpected stack value: Got '%v' but wanted '%v'", got, val)
			break
		}
	}
}

func assertStackLength(t *testing.T, s Stack, expected int) {
	if l := s.Len(); l != expected {
		t.Fatal("Bad stack size: Got %v but wanted %v", l, expected)
	}
}

func TestBasicStackFunctionality(t *testing.T) {
	s := NewStack()

	s.Push("Hello")
	assertStackLength(t, s, 1)
	s.Push("World")
	assertStackLength(t, s, 2)
	s.Push(42)
	assertStackLength(t, s, 3)

	checkStackPopOrder(t, s, 42, "World", "Hello")
}

func TestEmptyStack(t *testing.T) {
	s := NewStack()

	assertStackLength(t, s, 0)
	checkStackPopOrder(t, s, nil) // Expect nil
}

func TestReusingPoppedStack(t *testing.T) {
	s := NewStack()

	s.Push("Hello")
	s.Push("World")
	s.Push(42)

	// Pop all it's contents, expect nil for final (extra) value
	assertStackLength(t, s, 3)
	checkStackPopOrder(t, s, 42, "World", "Hello", nil)

	// Add more stuff ensuring it's still usable as a stack
	s.Push("Another value")
	s.Push("Final")
	assertStackLength(t, s, 2)
	checkStackPopOrder(t, s, "Final", "Another value", nil)
}
