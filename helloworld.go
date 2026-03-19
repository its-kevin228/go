package main

import {
	"fmt"
	"strconv"
}

var input1, input2, op string

func toFloat64(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func main() {

	fmt.Print("Enter first number: ")
	fmt.Scanln(&input1)

	fmt.Print("Enter second number: ")
	fmt.Scanln(&input2)

	fmt.Print("Enter operator (+, -, *, /): ")
	fmt.Scanln(&op)

	n1,err1:= toFloat64(input1)
	n2,err2:= toFloat64(input2)

	if err1 != nil {
		fmt.Println("Invalid input for first number.")
		return
	}
	if err2 != nil {
		fmt.Println("Invalid input for second number.")
		return
	}

	var result float64

	switch op {
	case "+":
		result = n1 + n2
	case "-":
		result = n1 - n2
	case "*":
		result = n1 * n2
	case "/":
		if n2 == 0 {
			fmt.Println("Error: Division by zero.")
			return
		}
		result = n1 / n2
	default:
		fmt.Println("Invalid operator.")
		return
	}

	fmt.Printf("Result: %.2f\n", result)
}