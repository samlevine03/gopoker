package benchmarks_test

import (
	"gopoker/poker"
	"testing"
	"time"
)

const (
	numHandsEval      = 10000
	numIterationsEval = 100
)

func BenchmarkHandEvaluation(b *testing.B) {
	b.Logf("Benchmarking hand evaluation speed...")

	evaluator := poker.NewEvaluator()
	var totalDuration time.Duration

	for i := 0; i < numIterationsEval; i++ {
		for j := 0; j < numHandsEval; j++ {
			deck := poker.NewDeck(int64(i * j))
			deck.Shuffle()

			cards := deck.Draw(7)

			start := time.Now()
			_ = evaluator.Evaluate(cards[:2], cards[2:])
			elapsed := time.Since(start)
			totalDuration += elapsed
		}
	}

	averageTime := totalDuration / (numHandsEval * numIterationsEval)
	b.Logf("\nHand Evaluation:\n  Total evaluations: %d\n  Total time: %v\n  Average time per hand: %v <---- \n  Hands per second: %.2f\n",
		numHandsEval*numIterationsEval, totalDuration, averageTime, float64(numHandsEval*numIterationsEval)/totalDuration.Seconds())
}
