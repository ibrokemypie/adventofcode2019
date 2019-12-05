package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please provide two variables (input url, input number)")
		os.Exit(1)
	}
	inputPath := os.Args[1]

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputNumber, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second arg must be a number")
		os.Exit(1)
	}

	fileString := string(content)
	intCodeStrings := strings.Split(strings.TrimSpace(fileString), ",")

	var intcode []int
	for i := 0; i < len(intCodeStrings); i++ {
		number, err := strconv.Atoi(intCodeStrings[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		intcode = append(intcode, number)
	}

	i := 0
	lastOutput := 0

	// run intcode
L:
	for i < len(intcode) {
		instruction := intcode[i]
		var opcode int

		// default to position mode (false)
		modeOne := 0
		modeTwo := 0

		instructionLen := int(math.Log10(float64(instruction))) + 1

		if instructionLen <= 2 && instructionLen != 0 {
			opcode = instruction
		} else if instructionLen <= 5 {
			opcode = instruction % 100
			modeOne = nthDigit(float64(instruction), 10, 3)
			modeTwo = nthDigit(float64(instruction), 10, 4)
		} else {
			fmt.Println("broke on instruction: " + strconv.Itoa(instruction) +
				" at position: " + strconv.Itoa(i))
			os.Exit(1)
		}

		if lastOutput != 0 && opcode != 99 {
			panic("incorrect output: " + strconv.Itoa(lastOutput))
		}

		switch opcode {
		case 1:
			{
				// add

				var valueOne int
				if modeOne == 0 {
					//position
					posOne := intcode[i+1]
					valueOne = intcode[posOne]
				} else if modeOne == 1 {
					// immediate
					valueOne = intcode[i+1]
				}

				var valueTwo int
				if modeTwo == 0 {
					//position
					posTwo := intcode[i+2]
					valueTwo = intcode[posTwo]
				} else if modeTwo == 1 {
					// immediate
					valueTwo = intcode[i+2]
				}

				outputPos := intcode[i+3]

				sum := valueOne + valueTwo
				intcode[outputPos] = sum
				i += 4
			}
		case 2:
			{
				// mult

				var valueOne int
				if modeOne == 0 {
					//position
					posOne := intcode[i+1]
					valueOne = intcode[posOne]
				} else if modeOne == 1 {
					// immediate
					valueOne = intcode[i+1]
				}

				var valueTwo int
				if modeTwo == 0 {
					//position
					posTwo := intcode[i+2]
					valueTwo = intcode[posTwo]
				} else if modeTwo == 1 {
					// immediate
					valueTwo = intcode[i+2]
				}

				outputPos := intcode[i+3]

				product := valueOne * valueTwo
				intcode[outputPos] = product
				i += 4
			}
		case 3:
			{
				// input
				inputPos := intcode[i+1]
				intcode[inputPos] = inputNumber
				i += 2
			}
		case 4:
			{
				//output

				var outputValue int
				if modeOne == 0 {
					//position
					outputPos := intcode[i+1]
					outputValue = intcode[outputPos]
				} else if modeOne == 1 {
					// immediate
					outputValue = intcode[i+1]
				}

				fmt.Println(outputValue)
				lastOutput = outputValue
				i += 2
			}
		case 99:
			{
				// break
				break L
			}
		default:
			{
				fmt.Println("invalid instruction: " + strconv.Itoa(instruction))
				os.Exit(1)
			}
		}
	}
}

func nthDigit(number float64, base float64, n float64) int {
	return int(number/math.Pow(base, n-1)) % int(base)
}
