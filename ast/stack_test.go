package ast

import (
	"reflect"
	"testing"
)

type stackValTest struct {
	Stack *Stack
	Want  *Stack
	Val   int
}

type stackTest struct {
	Stack *Stack
	Want  *Stack
}

func TestPush(t *testing.T) {
	for i, test := range []stackValTest{
		{&Stack{nil, 0, 0, 0}, &Stack{[]int{0}, 1, 0, 0}, 0},
		{&Stack{nil, 6, -3, -7}, &Stack{[]int{6}, 7, -3, -7}, 6},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{0, 1, 3}, 4, 0, 0}, 3},
	} {
		intEqual(t, i, test.Stack.Push(), test.Val)
		stackEqual(t, i, test.Stack, test.Want)
	}
}

func TestPop(t *testing.T) {
	for i, test := range []stackValTest{
		{&Stack{nil, 0, 0, 0}, &Stack{nil, 0, -1, -1}, -1},
		{&Stack{nil, 6, -3, -7}, &Stack{nil, 6, -4, -7}, -4},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{0}, 3, 0, 0}, 1},
	} {
		intEqual(t, i, test.Stack.Pop(), test.Val)
		stackEqual(t, i, test.Stack, test.Want)
	}
}

func TestPopN(t *testing.T) {
	for i, test := range []stackValTest{
		{&Stack{nil, 0, 0, 0}, &Stack{nil, 0, -1, -1}, 1},
		{&Stack{nil, 6, -3, -7}, &Stack{nil, 6, -4, -7}, 1},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{0}, 3, 0, 0}, 1},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{}, 3, -2, -2}, 4},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{0, 1}, 3, 0, 0}, 0},
	} {
		test.Stack.PopN(test.Val)
		stackEqual(t, i, test.Stack, test.Want)
	}

	checkPanic(t, -1, "stack: pop count must be positive: -1", func() {
		new(Stack).PopN(-1)
	})
}

func TestSwap(t *testing.T) {
	for i, test := range []stackTest{
		{&Stack{nil, 0, 0, 0}, &Stack{[]int{-1, -2}, 0, -2, -2}},
		{&Stack{nil, 6, -3, -7}, &Stack{[]int{-4, -5}, 6, -5, -7}},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{1, 0}, 3, 0, 0}},
		{&Stack{[]int{2}, 3, -1, -1}, &Stack{[]int{2, -2}, 3, -2, -2}},

		{&Stack{[]int{-1, -2}, 0, -2, -2}, &Stack{[]int{}, 0, 0, -2}},
		{&Stack{[]int{-4, -5}, 6, -5, -7}, &Stack{[]int{}, 6, -3, -7}},
		{&Stack{[]int{1, 0}, 3, 0, 0}, &Stack{[]int{0, 1}, 3, 0, 0}},
		{&Stack{[]int{2, -2}, 3, -2, -2}, &Stack{[]int{2}, 3, -1, -2}},
	} {
		test.Stack.Swap()
		stackEqual(t, i, test.Stack, test.Want)
	}
}

func TestSimplify(t *testing.T) {
	for i, test := range []stackTest{
		{&Stack{nil, 0, 0, 0}, &Stack{nil, 0, 0, 0}},
		{&Stack{[]int{0, 1}, 3, 0, 0}, &Stack{[]int{0, 1}, 3, 0, 0}},
		{&Stack{[]int{-1, -2}, 3, -1, -2}, &Stack{[]int{-2}, 3, 0, -2}},
		{&Stack{[]int{-3, -2, 0}, 2, -3, -3}, &Stack{[]int{0}, 2, -1, -3}},
	} {
		test.Stack.simplify()
		stackEqual(t, i, test.Stack, test.Want)
	}
}

func stackEqual(t *testing.T, testIndex int, got, want *Stack) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("test %d: got stack %v, want %v", testIndex, got, want)
	}
}

func intEqual(t *testing.T, testIndex int, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("test %d: got %d, want %d", testIndex, got, want)
	}
}

func checkPanic(t *testing.T, testIndex int, want interface{}, mightPanic func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r != want {
			t.Errorf("test %d: got panic %v, want panic %v", testIndex, r, want)
		}
	}()
	mightPanic()
}
