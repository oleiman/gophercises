// Package deck provides types and utilities for creating and processing decks of cards.
package deck

import (
	"math/rand"
	"sort"
	"time"
)

type Deck []Card
type DeckOption func(Deck) Deck

// NewDeck constructs a new Deck, configured by a variadic list of functional options.
func NewDeck(opts ...DeckOption) Deck {
	deck := make(Deck, 0, 52)
	for s := Spades; s <= Hearts; s++ {
		for v := Ace; v <= King; v++ {
			deck = append(deck, Card{Suit: s, Rank: v})
		}
	}

	for _, opt := range opts {
		deck = opt(deck)
	}
	return deck
}

// WithCustomSort produces a DeckOption to sort a Deck based on the provided less function.
func WithCustomSort(less func(d Deck) func(i, j int) bool) DeckOption {
	return func(d Deck) Deck {
		sort.Slice(d, less(d))
		return d
	}
}

// DefaultSort is a DeckOption for sorting a Deck by increasing absolute rank.
func DefaultSort(d Deck) Deck {
	sort.Slice(d, func(i, j int) bool {
		return d[i].Absolute() <= d[j].Absolute()
	})
	return d
}

// Shuffle is a DeckOption for randomly shuffling a Deck.
func Shuffle(d Deck) Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d),
		func(i, j int) {
			d[i], d[j] = d[j], d[i]
		})
	return d
}

// WithJokers produces a DeckOption to add a number of Jokers to the Deck.
func WithJokers(n int) DeckOption {
	return func(d Deck) Deck {
		for i := 0; i < n; i++ {
			d = append(d, Card{Suit: Joker})
		}
		return d
	}
}

// Filter produces a DeckOption to filter out those cards satisfying the
// provided predicate.
func Filter(exclude func(c Card) bool) DeckOption {
	return func(d Deck) Deck {
		result := make(Deck, 0, len(d))
		for _, c := range d {
			if exclude(c) {
				continue
			}
			result = append(result, c)
		}
		return result
	}
}

// Combine produces a Deckoption to append the provided Decks to the end
// of a Deck being constructed.
func Combine(others ...Deck) DeckOption {
	return func(d Deck) Deck {
		for _, other := range others {
			d = append(d, other...)
		}
		return d
	}
}
