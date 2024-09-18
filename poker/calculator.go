package poker

import (
	"runtime"
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

// CalculateEquity handles equity calculation for 2 to 9 players
func (c *Calculator) CalculateEquity(playerHands [][]uint32) []float64 {
	// Remove player hands from the deck
	for _, hand := range playerHands {
		c.deck.RemoveCards(hand)
	}
	remainingDeck := c.deck.cards
	possibleBoards := Combinations(remainingDeck, 5)

	numPlayers := len(playerHands)
	playerWins := make([]int, numPlayers)
	playerTies := make([]int, numPlayers)
	totalBoards := len(possibleBoards)

	// Use a channel to collect results
	results := make(chan [2][]int, totalBoards) // 2D array, first is wins, second is ties

	// Set up a worker pool with goroutines
	numWorkers := runtime.GOMAXPROCS(0)
	boardChunks := chunkBoards(possibleBoards, numWorkers)

	for _, boardChunk := range boardChunks {
		go func(boardChunk [][]uint32) {
			localWins := make([]int, numPlayers)
			localTies := make([]int, numPlayers)
			for _, board := range boardChunk {
				boardCards := []uint32{uint32(board[0]), uint32(board[1]), uint32(board[2]), uint32(board[3]), uint32(board[4])}
				scores := make([]int, numPlayers)
				bestScore := LookupTableMaxHighCard + 1

				// Evaluate each player's hand
				for i, hand := range playerHands {
					scores[i] = c.evaluator.Evaluate(hand, boardCards)
					if scores[i] < bestScore {
						bestScore = scores[i]
					}
				}

				// Count winners and check for ties
				var winners []int
				for i := 0; i < numPlayers; i++ {
					if scores[i] == bestScore {
						winners = append(winners, i)
					}
				}

				// If more than one player has the best score, it's a tie
				if len(winners) > 1 {
					for _, winner := range winners {
						localTies[winner]++
					}
				} else {
					localWins[winners[0]]++
				}
			}
			results <- [2][]int{localWins, localTies}
		}(boardChunk)
	}

	// Collect results from all workers
	for i := 0; i < numWorkers; i++ {
		res := <-results
		for j := 0; j < numPlayers; j++ {
			playerWins[j] += res[0][j]
			playerTies[j] += res[1][j]
		}
	}

	// Calculate equities
	playerEquities := make([]float64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		playerEquities[i] = (float64(playerWins[i]) + float64(playerTies[i])/float64(len(playerHands))) / float64(totalBoards) * 100
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
