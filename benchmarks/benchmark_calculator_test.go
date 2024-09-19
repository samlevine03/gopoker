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
		deck := poker.NewDeck(int64(i))
		deck.Shuffle()

		playerHands := [][]uint32{
			deck.Draw(2),
			deck.Draw(2),
		}

		calculator.SetDeck(deck)

		start := time.Now()
		_ = calculator.CalculateEquity(playerHands)
		elapsed := time.Since(start)
		totalDuration += elapsed
	}

	averageTime := totalDuration / numIterationsEquity
	b.Logf("\nEquity Calculation:\n  Total calculations: %d\n  Total time: %v\n  Average time per calculation: %v <---- \n  Calculations per second: %.2f\n",
		numIterationsEquity, totalDuration, averageTime, float64(numIterationsEquity)/totalDuration.Seconds())
}
