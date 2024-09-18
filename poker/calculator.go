package poker

import (
	"runtime"
	"sync"
)

type Calculator struct {
	deck      *Deck
	evaluator *Evaluator
}

func NewCalculator() *Calculator {
	return &Calculator{
		evaluator: NewEvaluator(),
	}
}

// SetDeck assigns a (shuffled?) deck to the calculator
func (c *Calculator) SetDeck(deck *Deck) {
	c.deck = deck
}

func (c *Calculator) CalculateEquity(playerHands [][]uint32) []float64 {
	// Remove player hands from the deck
	for _, hand := range playerHands {
		c.deck.RemoveCards(hand)
	}
	remainingDeck := c.deck.cards
	possibleBoards := Combinations(remainingDeck, 5)

	numPlayers := len(playerHands)
	totalBoards := len(possibleBoards)

	// Set up a worker pool with goroutines
	numWorkers := runtime.GOMAXPROCS(0)
	boardChunks := chunkBoards(possibleBoards, numWorkers)

	// Use a channel to collect results
	results := make(chan [2][]int, totalBoards) // 2D array, first is wins, second is ties

	var wg sync.WaitGroup
	wg.Add(len(boardChunks))

	for _, boardChunk := range boardChunks {
		go func(boardChunk [][]uint32) {
			defer wg.Done()
			localWins := make([]int, numPlayers)
			localTies := make([]int, numPlayers)
			scores := make([]int, numPlayers)

			for _, board := range boardChunk {
				bestScore := LookupTableMaxHighCard + 1
				winners := 0

				for i, hand := range playerHands {
					scores[i] = c.evaluator.Evaluate(hand, board)
					if scores[i] < bestScore {
						bestScore = scores[i]
						winners = 1
					} else if scores[i] == bestScore {
						winners++
					}
				}

				if winners > 1 {
					for i, score := range scores {
						if score == bestScore {
							localTies[i]++
						}
					}
				} else {
					for i, score := range scores {
						if score == bestScore {
							localWins[i]++
							break
						}
					}
				}
			}
			results <- [2][]int{localWins, localTies}
		}(boardChunk)
	}

	// Close the results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from all workers
	playerWins := make([]int, numPlayers)
	playerTies := make([]int, numPlayers)
	for res := range results {
		for j := 0; j < numPlayers; j++ {
			playerWins[j] += res[0][j]
			playerTies[j] += res[1][j]
		}
	}

	// Calculate equities
	playerEquities := make([]float64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		playerEquities[i] = (float64(playerWins[i]) + float64(playerTies[i])/float64(numPlayers)) / float64(totalBoards) * 100
	}

	return playerEquities
}

func chunkBoards(boards [][]uint32, numChunks int) [][][]uint32 {
	chunkSize := (len(boards) + numChunks - 1) / numChunks
	var chunks [][][]uint32
	for i := 0; i < len(boards); i += chunkSize {
		end := i + chunkSize
		if end > len(boards) {
			end = len(boards)
		}
		chunks = append(chunks, boards[i:end])
	}
	return chunks
}
