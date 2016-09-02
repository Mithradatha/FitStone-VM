package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// TODO: Add channels
func main() {

	arrays := map[uint32][]uint32{}
	var executionFinger uint32
	registers := [8]uint32{}
	fmt.Println(registers, "   -   ", arrays, executionFinger)

	if len(os.Args) > 1 {
		filename := os.Args[1]
		filePtr, err := os.Open(filename)
		if err == nil {
			defer filePtr.Close()
			buffer := make([]byte, 4)
			var programs []uint32
			for nBytes, err := filePtr.Read(buffer); nBytes != 0 && err == nil; nBytes, err = filePtr.Read(buffer) {
				//fmt.Printf("%b   %d   %d\n", buffer, binary.BigEndian.Uint32(buffer), binary.BigEndian.Uint32(buffer)>>28)
				programs = append(programs, binary.BigEndian.Uint32(buffer))
			}
			arrays[0] = programs

			// TODO: cycle()
			program := arrays[0][executionFinger]
			operator := program >> 28
			registerA := (program & 0x1FF) >> 6
			registerB := (program & 0x3F) >> 3
			registerC := (program & 0x7) // number & (2^nBits - 1)
			fmt.Printf("%d  %d, %d, %d\n", operator, registerA, registerB, registerC)
			fmt.Println(len(programs))
			executionFinger++
			switch operator {
			case 0: // conditionalMove
				if registers[registerC] != 0 {
					registers[registerA] = registers[registerB]
				}
			case 1: // arrayIndex
				registers[registerA] = arrays[registers[registerB]][registers[registerC]]
			case 2: // arrayAmendment
				arrays[registers[registerA]][registers[registerB]] = registers[registerC]
			case 3: // addition
				registers[registerA] = uint32(math.Mod(float64(registers[registerB]+registers[registerC]), math.Pow(2, 32)))
			case 4: // multiplication
				registers[registerA] = uint32(math.Mod(float64(registers[registerB]*registers[registerC]), math.Pow(2, 32)))
			case 5: // division
				if registers[registerC] != 0 {
					registers[registerA] = uint32(registers[registerB] / registers[registerC])
				}
			case 6: // notAnd
				registers[registerA] = ^(registers[registerB] & registers[registerC])
			case 7: // halt
				// TODO: Throw exception or alter isRunning bool to break cycle()
			case 8: // allocation

			case 9: // abandonment

			case 10: // output

			case 11: // input

			case 12: // loadProgram

			case 13: // orthography

			}
		}
	}
}
