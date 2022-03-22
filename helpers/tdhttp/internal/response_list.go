// Copyright (c) 2022, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package internal

import (
	"errors"
	"fmt"
	"sync"
)

type ResponseList struct {
	sync.Mutex
	responses map[string]*Response
	last      *Response
}

func NewResponseList() *ResponseList {
	return &ResponseList{
		responses: map[string]*Response{},
	}
}

func (rl *ResponseList) SetLast(resp *Response) {
	rl.Lock()
	defer rl.Unlock()

	rl.last = resp
}

func (rl *ResponseList) RecordLast(name string) error {
	rl.Lock()
	defer rl.Unlock()

	if rl.last == nil {
		return errors.New("no last response to record")
	}

	rl.last.Lock()
	defer rl.last.Unlock()

	if rl.last.name != "" {
		return fmt.Errorf("last response is already recorded as %q", rl.last.name)
	}

	rl.responses[name] = rl.last
	rl.last.name = name

	return nil
}

func (rl *ResponseList) Reset() {
	rl.Lock()
	defer rl.Unlock()

	for name := range rl.responses {
		delete(rl.responses, name)
	}
	rl.last = nil
}

func (rl *ResponseList) Get(name string) *Response {
	rl.Lock()
	defer rl.Unlock()
	return rl.responses[name]
}

func (rl *ResponseList) Last() *Response {
	rl.Lock()
	defer rl.Unlock()
	return rl.last
}
