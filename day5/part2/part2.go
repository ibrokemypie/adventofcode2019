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

		// fmt.Println(i)

		valueOne, valueTwo := getValues(i, intcode, modeOne, modeTwo)

		switch opcode {
		case 1:
			{
				// add

				outputPos := intcode[i+3]

				sum := valueOne + valueTwo
				intcode[outputPos] = sum
				i += 4
				// fmt.Println(strconv.Itoa(instruction) + ": add " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
			}
		case 2:
			{
				// mult

				outputPos := intcode[i+3]

				product := valueOne * valueTwo
				intcode[outputPos] = product
				i += 4
				// fmt.Println(strconv.Itoa(instruction) + ": multiply " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
			}
		case 3:
			{
				// input

				inputPos := intcode[i+1]
				intcode[inputPos] = inputNumber
				i += 2
				// fmt.Println(strconv.Itoa(instruction) + ": save " + strconv.Itoa(inputNumber) + " to position " + strconv.Itoa(inputPos))
			}
		case 4:
			{
				//output

				lastOutput = valueOne
				i += 2
				fmt.Println(strconv.Itoa(instruction) + ": output " + strconv.Itoa(valueOne))
			}
		case 5:
			{
				// jump if true (parameter, position)

				if valueOne != 0 {
					i = valueTwo
				} else {
					i += 3
				}
				// fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + " ==1 ")
			}
		case 6:
			{
				// jump if false

				if valueOne == 0 {
					i = valueTwo
				} else {
					i += 3
				}
				// fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + " ==0")
			}
		case 7:
			{
				// less than

				outputPos := intcode[i+3]

				if valueOne < valueTwo {
					intcode[outputPos] = 1
				} else {
					intcode[outputPos] = 0
				}
				i += 4

				// fmt.Println(strconv.Itoa(instruction) + ": set " + strconv.Itoa(outputPos) + " to 1 if " + strconv.Itoa(valueOne) + " smaller than " + strconv.Itoa(valueTwo))
			}
		case 8:
			{
				// equals

				outputPos := intcode[i+3]

				if valueOne == valueTwo {
					intcode[outputPos] = 1
				} else {
					intcode[outputPos] = 0
				}
				i += 4

				// fmt.Println(strconv.Itoa(instruction) + ": set " + strconv.Itoa(outputPos) + " to 1 if " + strconv.Itoa(valueOne) + " equal to " + strconv.Itoa(valueTwo))
			}
		case 99:
			{
				// break

				break L
			}
		default:
			{
				panic("invalid instruction: " + strconv.Itoa(instruction))
			}
		}
	}
}

func nthDigit(number float64, base float64, n float64) int {
	return int(number/math.Pow(base, n-1)) % int(base)
}

func getValues(pointer int, intcode []int, modeOne int, modeTwo int) (int, int) {
	var valueOne int
	var valueTwo int

	if pointer < len(intcode)-2 {
		if modeOne == 0 {
			//position
			posOne := intcode[pointer+1]
			valueOne = intcode[posOne]
		} else if modeOne == 1 {
			// immediate
			valueOne = intcode[pointer+1]
		}
	}

	if pointer < len(intcode)-2 {
		if modeTwo == 0 {
			//position
			posTwo := intcode[pointer+2]
			valueTwo = intcode[posTwo]
		} else if modeTwo == 1 {
			// immediate
			valueTwo = intcode[pointer+2]
		}
	}

	return valueOne, valueTwo
}
