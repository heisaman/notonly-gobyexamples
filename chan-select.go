package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func addStringToChan(chanN chan string) {
	for {
		chanN <- randStringRunes(10)
		time.Sleep(1 * time.Second)
	}

}

func main() {

	var chan1 = make(chan string, 10)
	var chan2 = make(chan string, 10)

	go addStringToChan(chan1)
	go addStringToChan(chan2)

	for {
		select {
		case e := <-chan1:
			fmt.Printf("Received random string from chan1: %s!\n", e)
		case e := <-chan2:
			fmt.Printf("Received random string from chan2: %s!\n", e)
		default:
			fmt.Println("No element in chan1 and chan2.")
			time.Sleep(1 * time.Second)
		}
	}
}
