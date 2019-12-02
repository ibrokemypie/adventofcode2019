package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide one variable (input url)")
		os.Exit(1)
	}
	inputPath := os.Args[1]

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	var intcode []int
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	// read comma separated opcodes into slice
	for scanner.Scan() {
		var data string
		for scanner.Text() != "," && scanner.Text() != "\n" {
			data += scanner.Text()
			scanner.Scan()
		}
		num, err := strconv.Atoi(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		intcode = append(intcode, num)
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
