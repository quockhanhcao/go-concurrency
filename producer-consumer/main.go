package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

// Producer is a type for structs that holds two channels: one for pizzas, with all
// information for a given order including whether it was made
// successfully or not, and another for handling end of processing (when we quit the channel)
type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

// PizzaOrder is a type for structs that holds all information for a pizza order
// It has the order number, a message indicating what happened to the order
// and a boolean indicating whether the order was successful or not
type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

// Close is a method for the Producer struct that closes the channel
// push something to the quit channel
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

// makePizza is a function that attemps to make a pizza
// generate random number from 1 to 12
// there are two failure cases.
// otherwise, it returns a pointer to a PizzaOrder struct
// each pizza takes a random amount of time to make
func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order number #%d\n", pizzaNumber)
		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d, it'll take %d seconds\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza #%d is ready", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	} else {
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
		}
	}
}

// pizzeria is a goroutine that runs in the background
// calls makePizza to make one pizza each time it iterates through the loop
// it executes until it receives something on the quit channel
// The quit channel does not receive anything until the consumer sends it
func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0
	// run forever or until we receive a quit notification
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				// close channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	// seed the random number generator
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	// print out message
	color.Cyan("The Pizzeria is open for business")
	color.Cyan("---------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer (one or more consumers)
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is ready", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The order failed")
			}
		} else {
			color.Cyan("Done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing the channel: %s", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("----------------------------------")
	color.Cyan("The Pizzeria is closed for business")
	color.Cyan("We made %d pizzas but failed to make %d pizzas with %d attemps in total", pizzasMade, pizzasFailed, total)
	switch {
	case pizzasFailed > 9:
		color.Red("Awful day")
	case pizzasFailed >= 6:
		color.Red("Bad day")
	case pizzasFailed >= 4:
		color.Yellow("An OK day")
	case pizzasFailed >= 2:
		color.Green("A good day")
	}
}
