package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	add = iota + 1
	mult
	input
	output
	jumptrue
	jumpfalse
	less
	equal
	adjbase
	end = 99
)

const (
	position = iota
	immediate
	relative
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
	intcode := parseIntcode(content)
	runIntcode(intcode, inputNumber)

}
func runIntcode(intcode []int, inputNumber int) {
	debug := false
	i := 0
	lastOutput := 0
	relativeBase := 0

	// run intcode
L:
	for i < len(intcode) {
		instruction := intcode[i]
		var opcode int

		// default to position mode (false)
		modeOne, modeTwo, modeThree, valueOne, valueTwo, posOne, posTwo, outputPos := 0, 0, 0, 0, 0, 0, 0, 0

		instructionLen := int(math.Log10(float64(instruction))) + 1

		if instructionLen <= 2 && instructionLen != 0 {
			opcode = instruction
		} else if instructionLen <= 5 {
			opcode = instruction % 100
			modeOne = nthDigit(float64(instruction), 10, 3)
			modeTwo = nthDigit(float64(instruction), 10, 4)
			modeThree = nthDigit(float64(instruction), 10, 5)
		} else {
			fmt.Println("broke on instruction: " + strconv.Itoa(instruction) +
				" at position: " + strconv.Itoa(i))
			os.Exit(1)
		}

		if lastOutput != 0 && opcode != 99 {
			panic("incorrect output: " + strconv.Itoa(lastOutput))
		}

		valueOne, valueTwo, posOne, posTwo, outputPos, intcode = getValues(i, intcode, opcode, modeOne, modeTwo, modeThree, relativeBase)

		if debug {
			fmt.Print("pointer: " + strconv.Itoa(i))
			fmt.Print(", instruction: " + strconv.Itoa(instruction))
			fmt.Print(", relbase: " + strconv.Itoa(relativeBase))
			fmt.Print(", pos1: " + strconv.Itoa(posOne))
			fmt.Print(", value1: " + strconv.Itoa(valueOne))
			fmt.Print(", pos2: " + strconv.Itoa(posTwo))
			fmt.Print(", value2: " + strconv.Itoa(valueTwo) + "\n")
		}

		switch opcode {
		case add:
			{
				// outputPos := intcode[i+3]

				if len(intcode) < outputPos {
					intcode = extendIntcode(intcode, outputPos)
				}

				sum := valueOne + valueTwo
				intcode[outputPos] = sum
				i += 4
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": add " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
				}
			}
		case mult:
			{
				// outputPos := intcode[i+3]

				if len(intcode) < outputPos {
					intcode = extendIntcode(intcode, outputPos)
				}

				product := valueOne * valueTwo
				intcode[outputPos] = product
				i += 4
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": multiply " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
				}
			}
		case input:
			{
				inputPos := posOne
				intcode[inputPos] = inputNumber
				i += 2
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": save " + strconv.Itoa(inputNumber) + " to position " + strconv.Itoa(inputPos))
				}
			}
		case output:
			{
				lastOutput = valueOne
				i += 2
				fmt.Println(strconv.Itoa(instruction) + ": output " + strconv.Itoa(valueOne))
			}
		case jumptrue:
			{
				if valueOne != 0 {
					i = valueTwo
				} else {
					i += 3
				}
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + " == 1 ")
				}
			}
		case jumpfalse:
			{
				if valueOne == 0 {
					i = valueTwo
				} else {
					i += 3
				}
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + "!= 1")
				}
			}
		case less:
			{
				// outputPos := intcode[i+3]

				if len(intcode) < outputPos {
					intcode = extendIntcode(intcode, outputPos)
				}

				if valueOne < valueTwo {
					intcode[outputPos] = 1
				} else {
					intcode[outputPos] = 0
				}
				i += 4
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": set " + strconv.Itoa(outputPos) + " to 1 if " + strconv.Itoa(valueOne) + " smaller than " + strconv.Itoa(valueTwo))
				}
			}
		case equal:
			{
				// outputPos := intcode[i+3]

				if len(intcode) < outputPos {
					intcode = extendIntcode(intcode, outputPos)
				}

				if valueOne == valueTwo {
					intcode[outputPos] = 1
				} else {
					intcode[outputPos] = 0
				}
				i += 4
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": set " + strconv.Itoa(outputPos) + " to 1 if " + strconv.Itoa(valueOne) + " equal to " + strconv.Itoa(valueTwo))
				}
			}
		case adjbase:
			{
				relativeBase += valueOne
				i += 2
				if debug {
					fmt.Println(strconv.Itoa(instruction) + ": adjust relativeBase by " + strconv.Itoa(valueOne))
				}
			}
		case end:
			{
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

func getValues(pointer int, intcode []int, opcode int, modeOne int, modeTwo int, modeThree int, relativeBase int) (int, int, int, int, int, []int) {
	var valueOne, valueTwo, posOne, posTwo, posThree int
	intcodeLength := len(intcode)

	if opcode != 99 {
		if pointer < len(intcode)-2 {
			if modeOne == position {
				//position
				posOne = intcode[pointer+1]
				if intcodeLength < posOne {
					intcode = extendIntcode(intcode, posOne)
				}
				valueOne = intcode[posOne]
			} else if modeOne == immediate {
				// immediate
				// posOne = intcode[pointer+1]
				valueOne = intcode[pointer+1]
			} else if modeOne == relative {
				posOne = relativeBase + intcode[pointer+1]
				if intcodeLength < posOne {
					intcode = extendIntcode(intcode, posOne)
				}
				valueOne = intcode[posOne]
			}

			if opcode != 3 && opcode != 4 {
				if modeTwo == position {
					//position
					posTwo = intcode[pointer+2]
					if intcodeLength < posTwo {
						intcode = extendIntcode(intcode, posTwo)
					}
					valueTwo = intcode[posTwo]
				} else if modeTwo == immediate {
					// immediate
					// posTwo = intcode[pointer+2]
					valueTwo = intcode[pointer+2]
				} else if modeTwo == relative {
					posTwo = relativeBase + intcode[pointer+2]
					if intcodeLength < posTwo {
						intcode = extendIntcode(intcode, posTwo)
					}
					valueTwo = intcode[posTwo]
				}

				if modeThree == position {
					//position
					posThree = intcode[pointer+3]
					if intcodeLength < posThree {
						intcode = extendIntcode(intcode, posThree)
					}
				} else if modeThree == immediate {
				} else if modeThree == relative {
					posThree = relativeBase + intcode[pointer+3]
					if intcodeLength < posThree {
						intcode = extendIntcode(intcode, posThree)
					}
				}
			}

		}
	}

	return valueOne, valueTwo, posOne, posTwo, posThree, intcode
}

func extendIntcode(intcode []int, value int) []int {
	newIntcode := make([]int, value+10, value+10)
	copy(newIntcode, intcode)
	return newIntcode
}

func parseIntcode(content []byte) []int {
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
	return intcode
}
