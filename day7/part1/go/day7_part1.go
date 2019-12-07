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
	end = 99
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

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	phaseSequence := [5]int{}
	intcode := parseIntcode(content)

	maxOutput := 0
	output := 0

	for five := 0; five < 5; five++ {
		for four := 0; four < 5; four++ {
			for three := 0; three < 5; three++ {
				for two := 0; two < 5; two++ {
				loop:
					for one := 0; one < 5; one++ {
						phaseSequence = [5]int{one, two, three, four, five}
						for k1, v1 := range phaseSequence {
							for k2, v2 := range phaseSequence {
								if k1 != k2 && v1 == v2 {
									continue loop
								}
							}
						}
						output = trySignal(phaseSequence, intcode, 0)
						if output > maxOutput {
							maxOutput = output
						}
					}
				}
			}
		}
	}
	fmt.Println(maxOutput)

}

func trySignal(phaseSequence [5]int, intcode []int, prevOutput int) int {
	output := prevOutput
	for _, phaseSignal := range phaseSequence {
		output = runIntcode(intcode, phaseSignal, output)
		if output != 0 {
			// fmt.Println(prevOutput)
		}
	}
	return output
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

func nthDigit(number float64, base float64, n float64) int {
	return int(number/math.Pow(base, n-1)) % int(base)
}

func getValues(pointer int, intcode []int, opcode int, modeOne int, modeTwo int) (int, int) {
	var valueOne int
	var valueTwo int

	if opcode != 99 {
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

		if opcode != 3 && opcode != 4 {
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
		}
	}

	return valueOne, valueTwo
}

func runIntcode(initialState []int, phaseSignal int, inputSignal int) int {
	i := 0
	lastOutput := 0
	intcode := make([]int, len(initialState))
	copy(intcode, initialState)
	inputCount := 0

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

		valueOne, valueTwo := getValues(i, intcode, opcode, modeOne, modeTwo)

		switch opcode {
		case add:
			{
				// add

				outputPos := intcode[i+3]

				sum := valueOne + valueTwo
				intcode[outputPos] = sum
				i += 4
				// fmt.Println(strconv.Itoa(instruction) + ": add " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
			}
		case mult:
			{
				// mult

				outputPos := intcode[i+3]

				product := valueOne * valueTwo
				intcode[outputPos] = product
				i += 4
				// fmt.Println(strconv.Itoa(instruction) + ": multiply " + strconv.Itoa(valueOne) + " and " + strconv.Itoa(valueTwo) + " output to position " + strconv.Itoa(outputPos))
			}
		case input:
			{
				// input

				inputPos := intcode[i+1]
				if inputCount == 0 {
					intcode[inputPos] = phaseSignal
				} else if inputCount == 1 {
					intcode[inputPos] = inputSignal
				}
				inputCount++
				i += 2
				// fmt.Println(strconv.Itoa(instruction) + ": save " + strconv.Itoa(inputNumber) + " to position " + strconv.Itoa(inputPos))
			}
		case output:
			{
				//output

				lastOutput = valueOne
				i += 2
				// fmt.Println(strconv.Itoa(instruction) + ": output " + strconv.Itoa(valueOne))
			}
		case jumptrue:
			{
				// jump if true (parameter, position)

				if valueOne != 0 {
					i = valueTwo
				} else {
					i += 3
				}
				// fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + " ==1 ")
			}
		case jumpfalse:
			{
				// jump if false

				if valueOne == 0 {
					i = valueTwo
				} else {
					i += 3
				}
				// fmt.Println(strconv.Itoa(instruction) + ": jump to " + strconv.Itoa(valueTwo) + " if " + strconv.Itoa(valueOne) + " ==0")
			}
		case less:
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
		case equal:
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
		case end:
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

	return lastOutput
}
