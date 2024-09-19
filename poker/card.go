package poker

import (
	"fmt"
)

const (
	StrRanks = "23456789TJQKA"
	StrSuits = "shdc"
	IntRanks = 13
)

var charRankToIntRank = map[rune]int{
	'2': 0, '3': 1, '4': 2, '5': 3, '6': 4, '7': 5,
	'8': 6, '9': 7, 'T': 8, 'J': 9, 'Q': 10, 'K': 11, 'A': 12,
}

var charSuitToIntSuit = map[rune]int{
	's': 1, 'h': 2, 'd': 4, 'c': 8,
}

var prettySuits = map[int]string{
	1: "♠", 2: "♥", 4: "♦", 8: "♣",
}

// NewCard creates a card from its string representation like "As" (Ace of Spades).
func NewCard(cardStr string) uint32 {
	rankInt := charRankToIntRank[rune(cardStr[0])]
	suitInt := charSuitToIntSuit[rune(cardStr[1])]

	bitRank := (1 << rankInt) << 16
	suit := suitInt << 12
	rank := rankInt << 8

	//        bitRank     suit rank   prime
	// +--------+--------+--------+--------+
	// |xxxbbbbb|bbbbbbbb|cdhsrrrr|xxpppppp|
	// +--------+--------+--------+--------+

	return uint32(bitRank | suit | rank | Primes[rankInt])
}

func GetRankInt(card uint32) int {
	return int((card >> 8) & 0xF)
}

func GetSuitInt(card uint32) int {
	return int((card >> 12) & 0xF)
}

func PrettyPrintCard(card uint32) {
	rankInt := GetRankInt(card)
	suitInt := GetSuitInt(card)
	fmt.Printf("[%s%s]", string(StrRanks[rankInt]), prettySuits[suitInt])
}

func PrintPrettyCards(cards []uint32) {
	for _, card := range cards {
		PrettyPrintCard(card)
	}
}
