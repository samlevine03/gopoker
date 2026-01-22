package main

import (
	"fmt"
	"gopoker/poker"
	"os"
	"strings"
)

func main() {
	// Check if we have command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go As,Ks Qh,Qd [more hands...]")
		fmt.Println("Example: go run main.go As,Ks Qh,Qd 6c,7c")
		os.Exit(1)
	}

	// Parse hands from command-line arguments
	var playerHands [][]uint32
	usedCards := make(map[string]bool)

	for i := 1; i < len(os.Args); i++ {
		handStr := os.Args[i]
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

	playerEquities := calculator.CalculateEquity(playerHands)

	for i := range playerHands {
		fmt.Printf("Player %d: %.2f%%\n", i+1, playerEquities[i])
	}
}
