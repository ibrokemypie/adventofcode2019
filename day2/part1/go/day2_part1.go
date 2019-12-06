package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide one variable (input url)")
		os.Exit(1)
	}
	inputPath := os.Args[1]

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var intcode []int
	fileString := string(content)
	intCodeStrings := strings.Split(strings.TrimSpace(fileString), ",")

	for i := 0; i < len(intCodeStrings); i++ {
		number, err := strconv.Atoi(intCodeStrings[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		intcode = append(intcode, number)
	}

	// apply challenge corrections
	intcode[1] = 12
	intcode[2] = 2

	// run intcode
	for i := 0; i < len(intcode)-1; i += 4 {
		opcode := intcode[i]
		pos1 := intcode[i+1]
		pos2 := intcode[i+2]
		outputPos := intcode[i+3]

		value1 := intcode[pos1]
		value2 := intcode[pos2]

		switch opcode {
		case 1:
			{
				sum := value1 + value2
				fmt.Println(strconv.Itoa(value1) + " + " + strconv.Itoa(value2) + " = " + strconv.Itoa(sum))
				intcode[outputPos] = sum
			}
		case 2:
			{
				product := value1 * value2
				fmt.Println(strconv.Itoa(value1) + " * " + strconv.Itoa(value2) + " = " + strconv.Itoa(product))
				intcode[outputPos] = product
			}
		}
	}
	fmt.Println("Value of pos 0: " + strconv.Itoa(intcode[0]))
}
