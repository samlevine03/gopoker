package poker

type LookupTable struct {
	flushLookup    map[int]int
	unsuitedLookup map[int]int
}

const (
	LookupTableMaxRoyalFlush    = 1
	LookupTableMaxStraightFlush = 10
	LookupTableMaxFourOfAKind   = 166
	LookupTableMaxFullHouse     = 322
	LookupTableMaxFlush         = 1599
	LookupTableMaxStraight      = 1609
	LookupTableMaxThreeOfAKind  = 2467
	LookupTableMaxTwoPair       = 3325
	LookupTableMaxPair          = 6185
	LookupTableMaxHighCard      = 7462
)

var LookupTableMaxToRankClass = map[int]int{
	LookupTableMaxRoyalFlush:    0,
	LookupTableMaxStraightFlush: 1,
	LookupTableMaxFourOfAKind:   2,
	LookupTableMaxFullHouse:     3,
	LookupTableMaxFlush:         4,
	LookupTableMaxStraight:      5,
	LookupTableMaxThreeOfAKind:  6,
	LookupTableMaxTwoPair:       7,
	LookupTableMaxPair:          8,
	LookupTableMaxHighCard:      9,
}

// NewLookupTable creates and returns a new lookup table
func NewLookupTable() *LookupTable {
	table := &LookupTable{
		flushLookup:    make(map[int]int),
		unsuitedLookup: make(map[int]int),
	}
	table.flushes()
	table.multiples()
	return table
}

// PrimeProductFromRankBits calculates the prime product using bitwise rank bits.
func PrimeProductFromRankBits(rankBits int) int {
	product := 1
	for i := 0; i < len(Primes); i++ {
		// If the ith bit is set
		if rankBits&(1<<i) != 0 {
			product *= Primes[i]
		}
	}
	return product
}

// PrimeProductFromHand calculates the prime product from a list of cards.
func PrimeProductFromHand(cards []uint32) int {
	product := 1
	for _, card := range cards {
		// Prime product is stored in the lower 6 bits of the card representation
		product *= int(card & 0x3F)
	}
	return product
}

// flushes generates the lookup for flushes and straight flushes
func (table *LookupTable) flushes() {
	// Straight flushes in rank order
	straightFlushes := []int{
		7936, 3968, 1984, 992, 496, 248, 124, 62, 31, 4111, // Royal flush down to 5-high straight flush
	}

	// Generate all other flushes
	var flushes []int
	gen := table.generateNextLexicographicalBitSequence(31) // 0b11111

	for i := 0; i < 1277+len(straightFlushes)-1; i++ {
		f := <-gen
		// Avoid adding straight flushes to the regular flushes
		notSF := true
		for _, sf := range straightFlushes {
			if f^sf == 0 {
				notSF = false
				break
			}
		}
		if notSF {
			flushes = append(flushes, f)
		}
	}

	// Reverse the flushes to rank from best to worst
	for i, j := 0, len(flushes)-1; i < j; i, j = i+1, j-1 {
		flushes[i], flushes[j] = flushes[j], flushes[i]
	}

	// Add straight flushes to flush lookup table
	rank := 1
	for _, sf := range straightFlushes {
		prime := PrimeProductFromRankBits(sf)
		table.flushLookup[prime] = rank
		rank++
	}

	// Add other flushes starting from max full house rank
	rank = LookupTableMaxFullHouse + 1
	for _, f := range flushes {
		prime := PrimeProductFromRankBits(f)
		table.flushLookup[prime] = rank
		rank++
	}

	// Generate straights and high cards
	table.straightAndHighCards(straightFlushes, flushes)
}

// straightAndHighCards generates lookup for straights and high cards
func (table *LookupTable) straightAndHighCards(straights []int, highCards []int) {
	rank := LookupTableMaxFlush + 1
	for _, s := range straights {
		prime := PrimeProductFromRankBits(s)
		table.unsuitedLookup[prime] = rank
		rank++
	}

	rank = LookupTableMaxPair + 1
	for _, h := range highCards {
		prime := PrimeProductFromRankBits(h)
		table.unsuitedLookup[prime] = rank
		rank++
	}
}

// multiples generates lookup for pairs, two pairs, three of a kind, full house, and four of a kind
func (table *LookupTable) multiples() {
	backwardsRanks := make([]int, len(Primes))
	for i := range Primes {
		backwardsRanks[i] = len(Primes) - 1 - i
	}

	// 1) Four of a kind
	rank := LookupTableMaxStraightFlush + 1
	for _, i := range backwardsRanks {
		kickers := make([]int, len(backwardsRanks))
		copy(kickers, backwardsRanks)
		for j := range kickers {
			if kickers[j] == i {
				continue
			}
			product := Primes[i] * Primes[i] * Primes[i] * Primes[i] * Primes[kickers[j]]
			table.unsuitedLookup[product] = rank
			rank++
		}
	}

	// 2) Full House
	rank = LookupTableMaxFourOfAKind + 1
	for _, i := range backwardsRanks {
		kickers := make([]int, len(backwardsRanks))
		copy(kickers, backwardsRanks)
		for j := range kickers {
			if kickers[j] == i {
				continue
			}
			product := Primes[i] * Primes[i] * Primes[i] * Primes[kickers[j]] * Primes[kickers[j]]
			table.unsuitedLookup[product] = rank
			rank++
		}
	}

	// 3) Three of a kind
	rank = LookupTableMaxStraight + 1
	for _, r := range backwardsRanks {
		kickers := make([]int, len(backwardsRanks))
		copy(kickers, backwardsRanks)
		for j := range kickers {
			if kickers[j] == r {
				continue
			}
			for k := j + 1; k < len(kickers); k++ {
				if kickers[k] == r {
					continue
				}
				product := Primes[r] * Primes[r] * Primes[r] * Primes[kickers[j]] * Primes[kickers[k]]
				table.unsuitedLookup[product] = rank
				rank++
			}
		}
	}

	// 4) Two Pair
	rank = LookupTableMaxThreeOfAKind + 1
	for i := range backwardsRanks {
		for j := i + 1; j < len(backwardsRanks); j++ {
			for k := 0; k < len(backwardsRanks); k++ {
				if k == i || k == j {
					continue
				}
				product := Primes[backwardsRanks[i]] * Primes[backwardsRanks[i]] * Primes[backwardsRanks[j]] * Primes[backwardsRanks[j]] * Primes[backwardsRanks[k]]
				table.unsuitedLookup[product] = rank
				rank++
			}
		}
	}

	// 5) One Pair
	rank = LookupTableMaxTwoPair + 1
	for i := range backwardsRanks {
		for j := 0; j < len(backwardsRanks); j++ {
			if i == j {
				continue
			}
			for k := j + 1; k < len(backwardsRanks); k++ {
				if k == i {
					continue
				}
				for l := k + 1; l < len(backwardsRanks); l++ {
					if l == i {
						continue
					}
					product := Primes[backwardsRanks[i]] * Primes[backwardsRanks[i]] * Primes[backwardsRanks[j]] * Primes[backwardsRanks[k]] * Primes[backwardsRanks[l]]
					table.unsuitedLookup[product] = rank
					rank++
				}
			}
		}
	}
}

// generateNextLexicographicalBitSequence generates the next lexicographical bit sequence
// For example, given the bit pattern 0b011 (3 in decimal), the next lexicographical bit sequence is 0b101 (5 in decimal).
// This function returns a channel that produces these sequences.
//
// Example:
// Input: 0b011 (3 in decimal)
// Output: 0b101 (5 in decimal), 0b110 (6 in decimal), ...
func (table *LookupTable) generateNextLexicographicalBitSequence(bits int) chan int {
	ch := make(chan int)
	go func() {
		for {
			t := (bits | (bits - 1)) + 1
			next := t | (((t & -t) / (bits & -bits)) >> 1) - 1
			ch <- next
			bits = next
		}
	}()
	return ch
}
