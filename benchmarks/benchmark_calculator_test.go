package benchmarks_test

import (
	"gopoker/poker"
	"testing"
	"time"
)

const (
	numIterationsEquity = 10
)

func BenchmarkCalculateEquity(b *testing.B) {
	b.Logf("Benchmarking equity calculation speed...")

	calculator := poker.NewCalculator()
	var totalDuration time.Duration

	for i := 0; i < numIterationsEquity; i++ {
		// Reset and shuffle the deck
		deck := poker.NewDeck(int64(i))
		deck.Shuffle()

		// Draw two random hands (2 cards each) from the deck
		playerHands := [][]uint32{
			deck.Draw(2), // Hand 1
			deck.Draw(2), // Hand 2
		}

		calculator.SetDeck(deck)

		// Start timer for this specific iteration
		start := time.Now()
		_ = calculator.CalculateEquity(playerHands)
		// Measure elapsed time
		elapsed := time.Since(start)
		totalDuration += elapsed
	}

	// Average and total time calculations
	averageTime := totalDuration / numIterationsEquity
	b.Logf("\nEquity Calculation:\n  Total calculations: %d\n  Total time: %v\n  Average time per calculation: %v <---- \n  Calculations per second: %.2f\n",
		numIterationsEquity, totalDuration, averageTime, float64(numIterationsEquity)/totalDuration.Seconds())
}
