# Poker Equity Calculator

## TODO:
 - [ ] CLI
 - [ ] Further optimization?
 - [ ] Benchmarks for hand evaluation speed
 - [ ] Usage instructions / examples (aka finish this README)

## Motivation

My goal was to write a pure Go implementation of a poker equity calculator.
I want to do so without using any external libraries, and without using any
external data sources. Furthermore, I didn't want to use a Monte Carlo simulation
to calculate the equity. 

### Hand Evaluation

My approach is based on [Cactus Kev's Poker Hand Evaluator.](http://suffe.cool/poker/evaluator.html)

If you don't want to click on any links, here's a summary:
> After enumerating and collapsing all the 2.6 million unique five-card poker hands, we wind up with just 7462 distinct poker hand values.

| Hand Value      | Unique  | Distinct |
|-----------------|---------|----------|
| Straight Flush  | 40      | 10       |
| 4 of a Kind     | 624     | 156      |
| Full House      | 3744    | 156      |
| Flush           | 5108    | 1277     |
| Straight        | 10200   | 10       |
| Three of a Kind | 54912   | 858      |
| Two Pair        | 123552  | 858      |
| One Pair        | 1098240 | 2860     |
| High Card       | 1302540 | 1277     |
| TOTAL           | 2598960 | 7462     |

We can turn [that](http://suffe.cool/poker/7462.html) into a lookup table. By representing each card with a prime number, we can get a unique product for any 5 card hand.

Cards can be represented by 32-bit integers as follows:

```
+--------+--------+--------+--------+
|xxxbbbbb|bbbbbbbb|cdhsrrrr|xxpppppp|
+--------+--------+--------+--------+

p = prime number of rank (deuce=2,trey=3,four=5,...,ace=41)
r = rank of card (deuce=0,trey=1,four=2,five=3,...,ace=12)
cdhs = suit of card (bit turned on based on suit of card)
b = bit turned on depending on rank of card
```

### Equity Calculation

Simple — we just enumerate through all possible run-outs and count the number of times each player wins/loses/ties!