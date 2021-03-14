// Golang program to illustrate the usage of 
// Sleep() function 

// Including main package 
package main

// Importing time and fmt 
import (
	"fmt"
	"time"
)

// Main function 
func main() {

	// Creating channel using 
	// make keyword 
	mychan1 := make(chan string, 2)

	// Calling Sleep function of go 
	go func() {
		time.Sleep(2 * time.Second)

		// Displayed after sleep overs 
		mychan1 <- "output1"
	}()

	// Select statement 
	select {

	// Case statement 
	case out := <-mychan1:
		fmt.Println(out)

	// Calling After method 
	case <-time.After(3 * time.Second):
		fmt.Println("timeout....1")
	}

	// Again Creating channel using 
	// make keyword 
	mychan2 := make(chan string, 2)

	// Calling Sleep method of go 
	go func() {
		time.Sleep(6 * time.Second)

		// Printed after sleep overs 
		mychan2 <- "output2"
	}()

	// Select statement 
	select {

	// Case statement 
	case out := <-mychan2:
		fmt.Println(out)

	// Calling After method 
	case <-time.After(3 * time.Second):
		fmt.Println("timeout....2")
	}
} 
