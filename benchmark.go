package main

import (
	"fmt"
	"gopoker/poker"
	"time"
)

const (
	numHandsEval        = 10000
	numIterationsEval   = 100
	numIterationsEquity = 10
)

func main() {

	fmt.Println("Benchmarking hand evaluation speed...")

	evaluator := poker.NewEvaluator()

	totalEvalDuration := time.Duration(0)

	for i := 0; i < numIterationsEval; i++ {
		for j := 0; j < numHandsEval; j++ {
			// Create and shuffle the deck
			deck := poker.NewDeck(int64(i * j))
			deck.Shuffle()

			// Draw 7 cards (2 for hand, 5 for board)
			cards := deck.Draw(7)

			start := time.Now()
			_ = evaluator.Evaluate(cards[:2], cards[2:])
			evalDuration := time.Since(start)

			totalEvalDuration += evalDuration
		}
	}

	averageEvalDuration := totalEvalDuration / (numHandsEval * numIterationsEval)
	fmt.Printf("\nTotal time for %d iterations of %d hands: %v\n", numIterationsEval, numHandsEval, totalEvalDuration)
	fmt.Printf("Average time per hand: %v\n", averageEvalDuration)
	fmt.Printf("Average hands per second: %.2f\n", 1/averageEvalDuration.Seconds())

	fmt.Println("\nBenchmarking equity calculation speed...")

	calculator := poker.NewCalculator()
	equityTotalDuration := time.Duration(0)

	for i := 0; i < numIterationsEquity; i++ {
		// Reset the deck and shuffle
		deck := poker.NewDeck(int64(i))
		deck.Shuffle()

		// Draw two random hands (2 cards each) from the deck
		playerHands := [][]uint32{
			deck.Draw(2), // Hand 1
			deck.Draw(2), // Hand 2
		}

		// Set the deck after drawing the player hands
		calculator.SetDeck(deck)

		// Time the equity calculation
		start := time.Now()
		_ = calculator.CalculateEquity(playerHands)
		equityDuration := time.Since(start)

		equityTotalDuration += equityDuration
	}

	averageEquityDuration := equityTotalDuration / numIterationsEquity
	fmt.Printf("\nTotal time for %d iterations of random equity calculation: %v\n", numIterationsEquity, equityTotalDuration)
	fmt.Printf("Average time per random equity calculation: %v\n", averageEquityDuration)
	fmt.Printf("Average random equity calculations per second: %.2f\n", 1/averageEquityDuration.Seconds())
}
