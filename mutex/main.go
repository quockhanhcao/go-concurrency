package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	// variable for bank balance
	var bankBalance int
    var balance sync.Mutex

	// print out starting value
	fmt.Printf("Initial bank balance: $%d.00\n", bankBalance)

	// define weekly revenue
	incomes := []Income{
		{
			Source: "Main Job", Amount: 1000,
		},
		{
			Source: "Gifts", Amount: 100,
		},
		{
			Source: "Investment", Amount: 500,
		},
	}

    wg.Add(len(incomes))

	// loop through 52 weeks and print out how much money is made, keep a running total
    for i, income := range incomes {
        go func(i int, income Income) {
            defer wg.Done()
            for week := 1; week <= 52; week++ {
                balance.Lock()
                temp := bankBalance
                temp += income.Amount
                bankBalance = temp
                balance.Unlock()

                fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)
            }
        }(i, income)
    }
    wg.Wait()
	// print out final balance
    fmt.Printf("Final bank balance: $%d.00\n", bankBalance)
}
