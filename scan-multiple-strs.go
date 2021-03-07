package main

import "fmt"

var rcCount = 0

func input(x []string, err error) []string {
	rcCount += 1
	fmt.Println("rc ", rcCount)
	if err != nil {
		return x
	}
	var d string
	n, err := fmt.Scanf("%s", &d)
	fmt.Println(n, err, d)
	if n == 1 {
		x = append(x, d)
	}
	return input(x, err)
}

func main() {
	fmt.Println("Enter input:")
	x := input([]string{}, nil)
	fmt.Println("Input:", x)
}
