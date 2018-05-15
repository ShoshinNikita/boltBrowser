package db

type offsetStack struct {
	offsets []int
}

func (stack *offsetStack) top() int {
	return stack.offsets[len(stack.offsets)-1]
}

// add 1 into end of slice
func (stack *offsetStack) add() {
	stack.offsets = append(stack.offsets, 1)
}

func (stack *offsetStack) del() {
	stack.offsets = stack.offsets[:len(stack.offsets)-1]
}

// lastElement++
func (stack *offsetStack) inc() {
	stack.offsets[len(stack.offsets)-1]++
}

// lastElement--
func (stack *offsetStack) dec() {
	stack.offsets[len(stack.offsets)-1]--
}
