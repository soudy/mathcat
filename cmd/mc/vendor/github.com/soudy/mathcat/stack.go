// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

type stack []interface{}

func (s *stack) Push(tok interface{}) {
	*s = append(*s, tok)
}

func (s stack) Top() interface{} {
	if s.Empty() {
		panic("top on empty stack")
	}
	return s[len(s)-1]
}

func (s *stack) Pop() interface{} {
	tok := s.Top()
	*s = (*s)[:len(*s)-1]

	return tok
}

func (s stack) Empty() bool {
	return len(s) == 0
}
