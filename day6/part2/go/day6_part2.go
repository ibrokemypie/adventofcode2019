package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	youMap := getParents(orbits, "YOU", []string{})
	santaMap := getParents(orbits, "SAN", []string{})

L:
	for youStep, youPar := range youMap {
		for santaStep, santaPar := range santaMap {
			if youPar == santaPar {
				steps := youStep + santaStep
				fmt.Println("Common node: " + youPar + " steps between: " + strconv.Itoa(steps))
				break L
			}
		}
	}

}

func getParents(orbits map[string]string, child string, parents []string) []string {
	parent := orbits[child]
	if parent != "COM" {
		parents = append(parents, parent)
		parents = getParents(orbits, parent, parents)
	}
	return parents
}
