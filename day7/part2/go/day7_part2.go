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

	intcode := parseIntcode(content)

	var maxOutput int
	opts := []int{5, 6, 7, 8, 9}
	allPermutations := permutations(opts)

	for _, phaseSequence := range allPermutations {
		newOutput := trySequence(phaseSequence, intcode)
		if newOutput > maxOutput {
			maxOutput = newOutput
		}
	}

	fmt.Println(maxOutput)

}

// https://stackoverflow.com/a/30226442
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func trySequence(phaseSequence []int, initialState []int) int {
	channels := [5]chan int{}
	for i := 0; i < 5; i++ {
		channels[i] = make(chan int)
		go runIntcode(initialState, channels[i])
		channels[i] <- phaseSequence[i]
	}

	lastOutput := 0

Top:
	for {
		for _, channel := range channels {
			select {
			case channel <- lastOutput:
			default:
				break Top
			}
			output := <-channel
			lastOutput = output
		}
	}
	return lastOutput
}

func runIntcode(initialState []int, inputChan chan int) {
	i := 0
	var lastOutput int
	intcode := make([]int, len(initialState))
	copy(intcode, initialState)

	var lastOpcode int

I:
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

		lastOpcode = opcode

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

				intcode[inputPos] = <-inputChan

				i += 2
				// fmt.Println(strconv.Itoa(instruction) + ": save " + strconv.Itoa(inputNumber) + " to position " + strconv.Itoa(inputPos))
			}
		case output:
			{
				//output

				lastOutput = valueOne
				inputChan <- lastOutput
				// return lastOutput, false
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
				break I
			}
		default:
			{
				panic("invalid instruction: " + strconv.Itoa(instruction))
			}
		}
	}
	inputChan <- 0
	inputChan <- lastOpcode
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
