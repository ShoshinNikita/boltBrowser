package db

// pagesStack keeps number of page of current and previous buckets
// Serves for returning on page, with which the user went before
type pagesStack struct {
	pages []int
}

func (stack *pagesStack) top() int {
	return stack.pages[len(stack.pages)-1]
}

// add 1 into end of slice
func (stack *pagesStack) add() {
	stack.pages = append(stack.pages, 1)
}

func (stack *pagesStack) del() {
	stack.pages = stack.pages[:len(stack.pages)-1]
}

// lastElement++
func (stack *pagesStack) inc() {
	stack.pages[len(stack.pages)-1]++
}

// lastElement--
func (stack *pagesStack) dec() {
	stack.pages[len(stack.pages)-1]--
}
