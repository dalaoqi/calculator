package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func Str2slice(message string) []int {
	numbersStr := strings.Fields(message)
	var numbers []int
	for _, numberStr := range numbersStr {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("Invalid input. Please enter only numbers.")
			continue
		}
		numbers = append(numbers, number)
	}
	return numbers
}
