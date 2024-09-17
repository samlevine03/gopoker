package main

import (
	"fmt"
	"gopoker/poker"
	"time"
)

func main() {
	playerHands := [][]uint32{
		{poker.NewCard("As"), poker.NewCard("Ks")},
		{poker.NewCard("Qh"), poker.NewCard("Qd")},
		{poker.NewCard("6c"), poker.NewCard("7c")},
	}

	deck := poker.NewDeck(0)
	calculator := poker.NewCalculator()
	calculator.SetDeck(deck)

	start := time.Now()
	playerEquities := calculator.CalculateEquity(playerHands)
	elapsed := time.Since(start)

	for i, hand := range playerHands {
		fmt.Printf("Player %d's Hand: ", i+1)
		poker.PrintPrettyCards(hand)
		fmt.Printf("\nPlayer %d's Preflop Equity: %.2f%%\n", i+1, playerEquities[i])
	}

	fmt.Printf("\nCalculation took %s\n", elapsed)
}
