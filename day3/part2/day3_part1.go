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

	fileString := string(content)

	lines := strings.Split(fileString, "\n")
	instructions := make([][]string, 2)
	instructions[0] = strings.Split(lines[0], ",")
	instructions[1] = strings.Split(lines[1], ",")

	var wires [][][2]int

	for _, line := range instructions {
		length := 0
		x, y := 0, 0
		var wire [][2]int

		for _, instruction := range line {
			direction := instruction[:1]
			distance, _ := strconv.Atoi(instruction[1:])
			length += distance

			for i := 0; i < distance; i++ {
				switch direction {
				case "U":
					{
						y++
					}
				case "D":
					{
						y--
					}
				case "R":
					{
						x++
					}
				case "L":
					{
						x--
					}
				}
				wire = append(wire, [2]int{x, y})
			}

		}

		wires = append(wires, wire)
	}

	var intercepts [][3]int

	for stepsOne, wireOne := range wires[0] {
		for stepsTwo, wireTwo := range wires[1] {
			if wireOne == wireTwo {
				intercepts = append(intercepts, [3]int{wireOne[0], wireOne[1], stepsOne + stepsTwo})
			}
		}

	}

	minSteps := float64(0)

	for _, intercept := range intercepts {

		if minSteps != 0 {
			minSteps = math.Min(minSteps, float64(intercept[2]))
		} else {
			minSteps = float64(intercept[2])
		}
	}

	// add two because step counts dont include the one step from 0,0 on either wire
	fmt.Println(minSteps + 2)
}
