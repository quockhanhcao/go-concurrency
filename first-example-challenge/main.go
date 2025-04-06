package main

import (
	"fmt"
	"sync"
)

var msg string

func updateMessage(s string) {
	msg = s
}

func printMessage(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(msg)
}

var wg sync.WaitGroup

func main() {

	// challenge: modify this code so that the calls to updateMessage() on lines
	// 28, 30, and 33 run as goroutines, and implement wait groups so that
	// the program runs properly, and prints out three different messages.
	// Then, write a test for all three functions in this program: updateMessage(),
	// printMessage(), and main().
	msg = "Hello, world!"
    msgArr := []string {
        "Hello, universe!",
        "Hello, cosmos!",
        "Hello, world!",
    }

	wg.Add(1)
	updateMessage(msgArr[0])
	go printMessage(&wg)
	wg.Wait()

	wg.Add(1)
	updateMessage(msgArr[1])
	go printMessage(&wg)
	wg.Wait()

	wg.Add(1)
	updateMessage(msgArr[2])
	go printMessage(&wg)
	wg.Wait()
}
