package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type layer struct {
	data  []int
	ones  int
	zeros int
	twos  int
}

func main() {
	width := 25
	height := 6

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

	var inputDigits []int
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		for scanner.Text() != "\n" {
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			inputDigits = append(inputDigits, num)
			scanner.Scan()
		}

	}

	layerCount := len(inputDigits) / (width * height)
	layers := []layer{}

	// 0: black 1: white 2: null
	image := make([]string, width*height+1)
	for k := range image {
		image[k] = ""
	}

	currentPos := 0

	for i := 0; i < layerCount; i++ {
		layer := layer{data: []int{}}
		posInLayer := 0
		// fmt.Println(image)
		for j := currentPos; j < ((i + 1) * width * height); j++ {
			layer.data = append(layer.data, inputDigits[j])
			currentPos = j + 1

			posInLayer++
			if image[posInLayer] == "" {
				switch inputDigits[j] {
				case 0:
					{
						image[posInLayer] = " "
					}

				case 1:
					{
						image[posInLayer] = "#"
					}

				case 2:
					{
						image[posInLayer] = ""
					}
				}
			}
		}
		layers = append(layers, layer)
	}

	for z := range image {
		fmt.Print(image[z])
		if z%width == 0 {
			fmt.Println()
		}
	}

	fmt.Println(len(image))
}
