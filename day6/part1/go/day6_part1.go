package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide one variable (input path)")
		os.Exit(1)
	}
	inputPath := os.Args[1]

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	orbits := map[string]string{}

	for scanner.Scan() {
		relation := strings.Split(scanner.Text(), ")")
		// orbits[parent] = child
		orbits[relation[1]] = relation[0]
	}

	distances := map[string]int{}
	for child := range orbits {
		distance := 0
		tempchild := child
		for tempchild != "COM" {
			distance++
			tempchild = orbits[tempchild]
		}
		distances[child] = distance
	}

	answer := 0
	for _, distance := range distances {
		answer += distance
	}
	fmt.Println(answer)

}
