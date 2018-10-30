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
	"log"
	"time"

	"encoding/binary"
	"github.com/tbrandon/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	serv.RegisterFunctionHandler(3, ReadRegisters)
	serv.RegisterFunctionHandler(4, ReadRegisters)
	serv.RegisterFunctionHandler(6, WriteRegisters)
	serv.RegisterFunctionHandler(16, WriteHoldingRegisters)

	err := serv.ListenTCP(":502")
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}

func ReadRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	// Check the request is within the allocated memory
	if endRegister > 65535 {
		return []byte{}, &mbserver.IllegalDataAddress
	}
	dataSize := numRegs // / 8
	if (numRegs % 8) != 0 {
		//dataSize++
	}
	dataSize = dataSize * 2
	data := make([]byte, 1+dataSize)
	data[0] = byte(dataSize)
	for i, value := range s.HoldingRegisters[register:endRegister] {
		// Return all 1s, regardless of the value in the DiscreteInputs array.
		value += 0xF00F
		if (i*2)+2 > dataSize {
			break
		}
		data[(i*2)+1] = byte(int(value / 256))
		data[(i*2)+2] = byte(value % 256)
	}
	fmt.Println("READ: ", data)
	return data, &mbserver.Success
}

func WriteRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, value := registerAddressAndValue(frame)
	// Check the request is within the allocated memory
	if register > 65535 {
		return []byte{}, &mbserver.IllegalDataAddress
	}

	data := make([]byte, 4)
	data[0] = byte(int(register / 256))
	data[1] = byte(register % 256)
	data[2] = byte(int(value / 256))
	data[3] = byte(value % 256)
	s.HoldingRegisters[register] = value
	fmt.Println("WRITE: ", data)

	return data, &mbserver.Success
}

func WriteHoldingRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numregs, _ := registerAddressAndNumber(frame)
	value := frame.GetData()[5:]
	// Check the request is within the allocated memory
	if register > 65535 {
		return []byte{}, &mbserver.IllegalDataAddress
	}

	if len(value)/2 != numregs {
		return []byte{}, &mbserver.IllegalDataAddress
	}

	for i := 0; i < len(value); i += 2 {
		b := value[i : i+2]
		fmt.Println("Data ", register+(i/2), b)
	}

	data := make([]byte, 4)
	data = frame.GetData()[0:4]

	fmt.Println("WRITE: ", data)

	return data, &mbserver.Success
}

func registerAddressAndNumber(frame mbserver.Framer) (register int, numRegs int, endRegister int) {
	data := frame.GetData()
	fmt.Println("DATA: ", data)
	register = int(binary.BigEndian.Uint16(data[0:2]))
	fmt.Println("REG: ", register)
	numRegs = int(binary.BigEndian.Uint16(data[2:4]))
	endRegister = register + numRegs
	return register, numRegs, endRegister
}

func registerAddressAndValue(frame mbserver.Framer) (int, uint16) {
	data := frame.GetData()
	fmt.Println("DATA: ", data)
	register := int(binary.BigEndian.Uint16(data[0:2]))
	fmt.Println("REG: ", register)
	value := binary.BigEndian.Uint16(data[2:4])
	return register, value
}
