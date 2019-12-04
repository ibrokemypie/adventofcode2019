package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	min := 153517
	max := 630395

	var count int

	for i := min; i <= max; i++ {
		digits := strings.Split((strconv.Itoa(i)), "")
		if checkDupe(digits) && checkAscending(digits) {
			fmt.Println(i)
			count++
		}
	}

	fmt.Println(count)
}

func checkDupe(digits []string) bool {
	for i := 0; i < len(digits)-1; i++ {
		count := 0
		j := i

		for (j < len(digits)-1) && (digits[j] == digits[j+1]) {
			count++
			j++
		}

		if count > 0 {
			// the count only looks forward, if it matches once that meants there were 2 dupes etc
			count++
			if count <= 2 {
				return true
			}
			i = j
		}
	}

	return false
}

func checkAscending(digits []string) bool {
	for i := 0; i < len(digits)-1; i++ {
		digOne, _ := strconv.Atoi(digits[i])
		digTwo, _ := strconv.Atoi(digits[i+1])
		if digOne > digTwo {
			return false
		}
	}

	return true
}
