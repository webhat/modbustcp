package main

import (
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

	err := serv.ListenTCP("127.0.0.1:502")
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
	dataSize := numRegs / 8
	if (numRegs % 8) != 0 {
		dataSize++
	}
	data := make([]byte, 1+dataSize)
	data[0] = byte(dataSize)
	for i := range s.DiscreteInputs[register:endRegister] {
		// Return all 1s, regardless of the value in the DiscreteInputs array.
		shift := uint(i) % 8
		data[1+i/8] |= byte(1 << shift)
	}
	return data, &mbserver.Success
}

func WriteRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	// Check the request is within the allocated memory
	if endRegister > 65535 {
		return []byte{}, &mbserver.IllegalDataAddress
	}
	dataSize := numRegs / 8
	if (numRegs % 8) != 0 {
		dataSize++
	}
	data := make([]byte, 1+dataSize)
	data[0] = byte(dataSize)
	for i := range s.DiscreteInputs[register:endRegister] {
		// Return all 1s, regardless of the value in the DiscreteInputs array.
		shift := uint(i) % 8
		data[1+i/8] |= byte(1 << shift)
	}
	return data, &mbserver.Success
}

func registerAddressAndNumber(frame mbserver.Framer) (register int, numRegs int, endRegister int) {
	data := frame.GetData()
	register = int(binary.BigEndian.Uint16(data[0:2]))
	numRegs = int(binary.BigEndian.Uint16(data[2:4]))
	endRegister = register + numRegs
	return register, numRegs, endRegister
}
