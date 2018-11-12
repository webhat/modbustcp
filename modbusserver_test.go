/**
MIT License

Copyright (c) 2018 Daniel W. Crompton

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"testing"

	"github.com/tbrandon/mbserver"
)

type Framer interface {
	mbserver.Framer
}

type myFrame struct {
	data []byte
}

func (f *myFrame) Bytes() []byte                              { return f.data }
func (f *myFrame) GetData() []byte                            { return f.data }
func (f *myFrame) GetFunction() uint8                         { return uint8(0) }
func (f *myFrame) SetException(exception *mbserver.Exception) { return }
func (f *myFrame) SetData(data []byte)                        { f.data = data }
func (f *myFrame) Copy() mbserver.Framer {
	copy := *f
	return &copy
}

func TestSingleRead(t *testing.T) {
	s := mbserver.NewServer()
	f := myFrame{data: []byte{0xDE, 0xAD, 0x00, 0x01}}

	b, err := ReadRegisters(s, mbserver.Framer(&f))

	fmt.Println(b)
	if len(b) <= 0 {
		t.Errorf("Error: %v", err)
	}
}

func TestMultipleRead(t *testing.T) {
	s := mbserver.NewServer()
	f := myFrame{data: []byte{0xDE, 0xAD, 0x00, 50}}

	b, err := ReadRegisters(s, mbserver.Framer(&f))

	fmt.Println(b)
	if len(b) <= 0 {
		t.Errorf("Error: %v", err)
	}
}
