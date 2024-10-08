package poker

import (
	"math/rand"
	"time"
)

type Deck struct {
	cards []uint32
}

// fullDeck is the static full deck, initialized once
var fullDeck []uint32
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewDeck(seed int64) *Deck {
	if seed != 0 {
		rng = rand.New(rand.NewSource(seed))
	}
	return &Deck{
		cards: GetFullDeck(),
	}
}

func (d *Deck) Shuffle() {
	rng.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw n cards from the deck
func (d *Deck) Draw(n int) []uint32 {
	if n > len(d.cards) {
		n = len(d.cards)
	}
	drawn := d.cards[len(d.cards)-n:]  // Slice out the last n cards
	d.cards = d.cards[:len(d.cards)-n] // Shrink the deck
	return drawn
}

// GetFullDeck generates the full 52-card deck
func GetFullDeck() []uint32 {
	if len(fullDeck) > 0 {
		return append([]uint32(nil), fullDeck...)
	}

	for _, rank := range StrRanks {
		for _, suit := range StrSuits {
			fullDeck = append(fullDeck, NewCard(string(rank)+string(suit)))
		}
	}

	return append([]uint32(nil), fullDeck...)
}

// RemoveCards removes specific cards from the deck.
func (d *Deck) RemoveCards(cardsToRemove []uint32) {
	removeMap := make(map[uint32]bool)
	for _, card := range cardsToRemove {
		removeMap[card] = true
	}

	newSize := 0
	for _, card := range d.cards {
		if !removeMap[card] {
			d.cards[newSize] = card
			newSize++
		}
	}

	d.cards = d.cards[:newSize]
}
