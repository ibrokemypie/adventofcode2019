package main

import (
	"bufio"
	"fmt"
	"math"
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

	minZeros := math.MaxInt32
	layerMinZeros := layer{}

	currentPos := 0
	for i := 0; i < layerCount; i++ {
		layer := layer{data: []int{}}
		for j := currentPos; j < ((i + 1) * width * height); j++ {
			layer.data = append(layer.data, inputDigits[j])
			currentPos = j + 1
			switch inputDigits[j] {
			case 0:
				{
					layer.zeros++
				}
			case 1:
				{
					layer.ones++
				}
			case 2:
				{
					layer.twos++
				}
			default:
				{
				}
			}
		}
		layers = append(layers, layer)
		if layer.zeros < minZeros {
			minZeros = layer.zeros
			layerMinZeros = layer
		}
	}

	fmt.Println(layerMinZeros.ones * layerMinZeros.twos)
}
