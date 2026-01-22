package main

import (
	"flag"
	"fmt"
	"gopoker/poker"
	"os"
	"strings"
)

func main() {
	// Define -board flag for known board cards
	boardFlag := flag.String("board", "", "Known board cards (e.g., 'Jh,Tc,9d' for flop or 'Jh,Tc,9d,8s' for flop+turn)")
	flag.Parse()

	args := flag.Args()

	// Check if we have command-line arguments
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go [options] As,Ks Qh,Qd [more hands...]")
		fmt.Println("Options:")
		fmt.Println("  -board string  Known board cards (flop, flop+turn, or full board)")
		fmt.Println("Example: go run main.go -board Jh,Tc,9d As,Ks Qh,Qd")
		os.Exit(1)
	}

	// Parse board cards if provided
	var knownBoard []uint32
	usedCards := make(map[string]bool)

	if *boardFlag != "" {
		boardCardsStr := strings.Split(*boardFlag, ",")
		if len(boardCardsStr) < 3 || len(boardCardsStr) > 5 {
			fmt.Printf("Error: Board must have 3-5 cards (flop, flop+turn, or full board), got %d\n", len(boardCardsStr))
			os.Exit(1)
		}

		knownBoard = make([]uint32, len(boardCardsStr))
		for i, cardStr := range boardCardsStr {
			cardStr = strings.TrimSpace(cardStr)
			if usedCards[cardStr] {
				fmt.Printf("Error: Duplicate card in board: %s\n", cardStr)
				os.Exit(1)
			}
			usedCards[cardStr] = true
			knownBoard[i] = poker.NewCard(cardStr)
		}
	}

	// Parse hands from command-line arguments
	var playerHands [][]uint32

	for i := 0; i < len(args); i++ {
		handStr := args[i]
		cards := strings.Split(handStr, ",")

		if len(cards) != 2 {
			fmt.Printf("Error: Each hand must have exactly 2 cards, got %d in '%s'\n", len(cards), handStr)
			os.Exit(1)
		}

		hand := make([]uint32, 2)
		for j, cardStr := range cards {
			cardStr = strings.TrimSpace(cardStr)

			// Check for duplicates
			if usedCards[cardStr] {
				fmt.Printf("Error: Duplicate card: %s\n", cardStr)
				os.Exit(1)
			}
			usedCards[cardStr] = true

			hand[j] = poker.NewCard(cardStr)
		}

		playerHands = append(playerHands, hand)
	}

	if len(playerHands) < 2 || len(playerHands) > 9 {
		fmt.Printf("Error: Must have 2-9 players, got %d\n", len(playerHands))
		os.Exit(1)
	}

	deck := poker.NewDeck(0)
	calculator := poker.NewCalculator()
	calculator.SetDeck(deck)

	var playerEquities []float64
	if len(knownBoard) > 0 {
		playerEquities = calculator.CalculateEquityWithBoard(playerHands, knownBoard)
	} else {
		playerEquities = calculator.CalculateEquity(playerHands)
	}

	for i := range playerHands {
		fmt.Printf("Player %d: %.2f%%\n", i+1, playerEquities[i])
	}
}
