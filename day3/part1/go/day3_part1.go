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

	var intercepts [][2]int

	for _, wireOne := range wires[0] {
		for _, wireTwo := range wires[1] {
			if wireOne == wireTwo {
				intercepts = append(intercepts, wireOne)
			}
		}

	}

	minDistance := float64(0)

	for _, intercept := range intercepts {
		var distance float64

		distance += math.Abs(float64(intercept[0]))
		distance += math.Abs(float64(intercept[1]))

		if minDistance != 0 {
			minDistance = math.Min(minDistance, distance)
		} else {
			minDistance = distance
		}
	}

	fmt.Println(minDistance)
}
