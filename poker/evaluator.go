package poker

type Evaluator struct {
	table        *LookupTable
	HandLength   int
	BoardLength  int
	handSizeEval map[int]func([]uint32) int
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		table:       NewLookupTable(),
		HandLength:  2,
		BoardLength: 5,
	}
	e.handSizeEval = map[int]func([]uint32) int{
		5: e._five,
		7: e._seven,
	}
	return e
}

// Evaluate calculates the strength of the hand
func (e *Evaluator) Evaluate(hand []uint32, board []uint32) int {
	allCards := append(hand, board...)
	return e.handSizeEval[len(allCards)](allCards)
}

// _five evaluates a 5-card hand
func (e *Evaluator) _five(cards []uint32) int {
	flush := cards[0] & cards[1] & cards[2] & cards[3] & cards[4] & 0xF000

	if flush != 0 {
		handOR := (cards[0] | cards[1] | cards[2] | cards[3] | cards[4]) >> 16
		prime := PrimeProductFromRankBits(int(handOR))
		return e.table.flushLookup[prime]
	}
	prime := PrimeProductFromHand(cards)
	return e.table.unsuitedLookup[prime]
}

// _seven evaluates all 7 choose 5 combinations of a 7-card hand
func (e *Evaluator) _seven(cards []uint32) int {
	// Check for flush possibility
	suitCounts := [4]int{}
	for _, card := range cards {
		suit := (card >> 12) & 0xF
		if suit < 4 { // Ensure suit is valid
			suitCounts[suit]++
		}
	}

	for suit, count := range suitCounts {
		if count >= 5 {
			flushCards := make([]uint32, 0, 7)
			for _, card := range cards {
				if int(card>>12&0xF) == suit {
					flushCards = append(flushCards, card)
				}
			}
			return e._flushSeven(flushCards)
		}
	}

	minScore := LookupTableMaxHighCard
	combinations := Combinations(cards, 5)

	for _, combo := range combinations {
		score := e._five(combo)
		if score < minScore {
			minScore = score
		}
	}
	return minScore
}

func (e *Evaluator) _flushSeven(flushCards []uint32) int {
	minScore := LookupTableMaxHighCard
	combinations := Combinations(flushCards, 5)
	for _, combo := range combinations {
		score := e._five(combo)
		if score < minScore {
			minScore = score
		}
	}
	return minScore
}

// GetRankClass returns the rank class of a hand
func (e *Evaluator) GetRankClass(handRank int) int {
	switch {
	case handRank >= 0 && handRank <= LookupTableMaxRoyalFlush:
		return LookupTableMaxToRankClass[LookupTableMaxRoyalFlush]
	case handRank <= LookupTableMaxStraightFlush:
		return LookupTableMaxToRankClass[LookupTableMaxStraightFlush]
	case handRank <= LookupTableMaxFourOfAKind:
		return LookupTableMaxToRankClass[LookupTableMaxFourOfAKind]
	case handRank <= LookupTableMaxFullHouse:
		return LookupTableMaxToRankClass[LookupTableMaxFullHouse]
	case handRank <= LookupTableMaxFlush:
		return LookupTableMaxToRankClass[LookupTableMaxFlush]
	case handRank <= LookupTableMaxStraight:
		return LookupTableMaxToRankClass[LookupTableMaxStraight]
	case handRank <= LookupTableMaxThreeOfAKind:
		return LookupTableMaxToRankClass[LookupTableMaxThreeOfAKind]
	case handRank <= LookupTableMaxTwoPair:
		return LookupTableMaxToRankClass[LookupTableMaxTwoPair]
	case handRank <= LookupTableMaxPair:
		return LookupTableMaxToRankClass[LookupTableMaxPair]
	case handRank <= LookupTableMaxHighCard:
		return LookupTableMaxToRankClass[LookupTableMaxHighCard]
	default:
		panic("Invalid hand rank")
	}
}
