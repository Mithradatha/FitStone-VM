package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"bufio"
	"errors"
	"strconv"
	"time"
)

const maxUint32 = 4294967295
const maxUint8 = 255
const increment = 256
const stackSize = 562
const registerSize = 8

type AddressStack struct {
    unallocatedMemory chan uint32
    maxVal chan uint32
}

func (o *AddressStack) Init() {
    o.unallocatedMemory = make(chan uint32, stackSize)
    o.maxVal = make(chan uint32, 1)
    
    for i, j := uint32(0), uint32(increment); i < j; i++ {
        o.unallocatedMemory <- i
    }

    o.maxVal <- increment
}

func (o *AddressStack) Print() {
    
    close(o.unallocatedMemory)
    close(o.maxVal)

    fmt.Println("Largest Address: ", <- o.maxVal)
    fmt.Println("\nUnallocated Addresses:")

    for elem := range o.unallocatedMemory {
        fmt.Println(elem)
    }
}

func (o *AddressStack) populate() error {
    i := <- o.maxVal
    j := i + increment
    
    if j > maxUint32 {
        return errors.New("Memory Full")
    }

    for ; i < j; i++ {
        o.unallocatedMemory <- i
    }
    o.maxVal <- j
    return nil
}

func (o *AddressStack) Push(elem uint32) {
    if len(o.unallocatedMemory) == stackSize {
        o.Pop()
    }
    o.unallocatedMemory <- elem
}

func (o *AddressStack) Pop() (elem uint32, err error) {
    elem = <- o.unallocatedMemory
    if len(o.unallocatedMemory) < 50 {
        err = o.populate()
    }
    return
}

func main() {
	start := time.Now()
	output, _ := os.Create("sandmark-output.txt")

	var executionFinger uint32

	addresses := new(AddressStack)
	arrays := map[uint32][]uint32{}
	registers := [registerSize]uint32{}
	
	addresses.Init()

	if len(os.Args) > 1 {
		
		filename := os.Args[1]
		filePtr, err := os.Open(filename)
		
		if err == nil {
			
			defer filePtr.Close()
			
			buffer := make([]byte, 4)
			var programs []uint32
			
			for nBytes, err := filePtr.Read(buffer); nBytes != 0 && err == nil; nBytes, err = filePtr.Read(buffer) {
				programs = append(programs, binary.BigEndian.Uint32(buffer))
			}

			index0, _ := addresses.Pop()
			arrays[index0] = programs

			// cycle
			for isRunning := true; isRunning; {
				
				program := arrays[0][executionFinger]
				operator := program >> 28
				
				registerA := (program & 0x1FF) >> 6
				registerB := (program & 0x3F) >> 3
				registerC := (program & 0x7) // number & (2^nBits - 1)
				
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
					registers[registerA] = uint32(registers[registerB] / registers[registerC])
				case 6: // notAnd
					registers[registerA] = ^(registers[registerB] & registers[registerC])
				case 7: // halt
					isRunning = false
				case 8: // allocation
					address, _ := addresses.Pop()
					arrays[address] = make([]uint32, registers[registerC])
					registers[registerB] = address
				case 9: // abandonment
					address := registers[registerC]
					delete(arrays, address)
					addresses.Push(address)
				case 10: // output
					fmt.Fprintf(output, "%v", string(registers[registerC]))
				case 11: // input
					input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
					intStr, _ := strconv.ParseUint(input, 10, 32)
					registers[registerC] = uint32(intStr)
				case 12: // loadProgram
					if registers[registerB] != 0 {
						arr := arrays[registers[registerB]]
						tmp := make([]uint32, len(arr))
						copy(tmp, arr)
						arrays[0] = tmp
					}
					executionFinger = registers[registerC]
				case 13: // orthography
					registers[((program >> 25) & 0x7)] = (program & 0x1FFFFFF)
				}
			}	
		}
	}
	fmt.Println("Start: ", start)
	fmt.Println("End: ", time.Now())
	fmt.Println("Duration: ", time.Since(start))
}
