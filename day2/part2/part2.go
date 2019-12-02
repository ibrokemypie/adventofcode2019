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

	var initialState []int
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
		initialState = append(initialState, num)
	}

	for n := 0; n <= 99; n++ {
		for v := 0; v <= 99; v++ {
			intcode := make([]int, len(initialState))
			copy(intcode, initialState)

			// apply challenge corrections
			intcode[1] = n
			intcode[2] = v

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
						intcode[outputPos] = sum
					}
				case 2:
					{
						product := value1 * value2
						intcode[outputPos] = product
					}

				}
			}

			if intcode[0] == 19690720 {
				answer := 100*n + v
				fmt.Println("noun: " + strconv.Itoa(n) + ", verb: " + strconv.Itoa(v))
				fmt.Println("100 * " + strconv.Itoa(n) + " + " + strconv.Itoa(v) + " = " + strconv.Itoa(answer))
			}
		}

	}

}
