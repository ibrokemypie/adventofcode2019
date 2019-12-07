package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinate struct {
	x int
	y int
}

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

	var intersects []coordinate = getOverlappedPositions(getVisitied(instructions[0]), getVisitied(instructions[1]))
	var distances []int = []int{}

	for _, point := range intersects {
		distance := math.Abs(float64(0-point.x)) + math.Abs(float64(0-point.y))
		distances = append(distances, int(distance))
	}

	var minDistance int
	for i, e := range distances {
		if i == 0 || e < minDistance {
			minDistance = e
		}
	}

	fmt.Println(minDistance)
}

func getOverlappedPositions(path1 map[string]coordinate, path2 map[string]coordinate) []coordinate {
	var overlaps []coordinate = []coordinate{}

	for key, value := range path1 {
		if _, ok := path2[key]; ok {
			overlaps = append(overlaps, value)
		}
	}

	return overlaps
}

func getVisitied(line []string) map[string]coordinate {
	var position coordinate
	visited := make(map[string]coordinate)

	for _, instruction := range line {
		direction := instruction[:1]
		distance, _ := strconv.Atoi(instruction[1:])

		for i := 0; i < distance; i++ {
			switch direction {
			case "U":
				{
					position.y++
				}
			case "D":
				{
					position.y--
				}
			case "R":
				{
					position.x++
				}
			case "L":
				{
					position.x--
				}
			}
			visited[strconv.Itoa(position.x)+" "+strconv.Itoa(position.y)] = position
		}

	}
	return visited
}
