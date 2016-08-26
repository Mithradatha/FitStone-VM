package main

import "os"
import "fmt"

type Platter struct {
	value uint32
}

type Register struct {
	platter Platter
}

func main() {

	var inputFileName string

	if len(os.Args) > 1 {
		inputFileName = os.Args[1]
	}

	outputFileName := "sandmark-output.txt"
	fmt.Println(inputFileName)
	fmt.Println(outputFileName)

}
