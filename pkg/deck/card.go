package deck

import (
	"fmt"
)

//go:generate stringer -type=Suit
//go:generate stringer -type=Rank

type Suit int

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker
)

// Type Rank correspon
type Rank int

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Type Card provides the elements of composite type Deck.
// Each card encapsulates the card's Suit and its Rank.
type Card struct {
	Suit
	Rank
}

// Absolute returns the absolute rank of the card.
// Each card in a standard deck should have a unique Absolute rank.
func (c Card) Absolute() int {
	if c.Suit == Joker {
		return 53
	}
	return int(c.Suit)*(int(King)) + int(c.Rank)
}

// String provides a friendly string representation for the card.
// e.g. Card{Rank: Ace, Suit: Spades}.String() -> "Ace of Spades"
func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}
