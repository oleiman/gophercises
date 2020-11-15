package deck

import (
	"testing"
)

func TestCreateDeck(t *testing.T) {
	d := NewDeck()
	if len(d) != 52 {
		t.Errorf("Created a deck with %d cards...", len(d))
	}
}

var absoluteTestCases = [][]interface{}{
	{Card{Spades, Ace}, 1},
	{Card{Diamonds, Ace}, 14},
	{Card{Hearts, King}, 52},
}

func TestAbsolute(t *testing.T) {
	for _, c := range absoluteTestCases {
		r := c[0].(Card).Absolute()
		if r != c[1] {
			t.Errorf("%s  should have abslute rank %d (got %d)", c[0], c[1], r)
		}
	}
}

func TestCreateSortedDeck(t *testing.T) {
	reverse := func(d Deck) func(i, j int) bool {
		return func(i, j int) bool {
			return d[i].Absolute() > d[j].Absolute()
		}
	}
	d := NewDeck(WithCustomSort(reverse))
	expected := Card{Hearts, King}
	if d[0].Absolute() != expected.Absolute() {
		t.Errorf("Wrong order for reverse sort...first card was %s expected %s",
			d[0], expected)
	}

	d = NewDeck(DefaultSort)
	expected = Card{Spades, Ace}
	if d[0].Absolute() != expected.Absolute() {
		t.Errorf("Wrong order for default sort... first card was %s expected %s",
			d[0], expected)
	}
}

func TestShuffle(t *testing.T) {
	d := NewDeck(DefaultSort)
	dShuf := NewDeck(Shuffle)

	diff := 0
	for i, v := range dShuf {
		if d[i] != v {
			diff++
		}
	}
	// this is an incredibly dumb test that will eventually fail
	if diff < 10 {
		t.Errorf("Not enough entropy in shuffled deck (diff=%d)", diff)
	}
}

func TestWithJokers(t *testing.T) {
	d := NewDeck(WithJokers(2), DefaultSort)
	if len(d) < 54 {
		t.Errorf("Only %d jokers were added", len(d)-52)
	}
	if d[52].Suit != Joker {
		t.Errorf("Didn't add a joker, instead got %s", d[52])
	}
}

func TestFilter(t *testing.T) {
	exclude := func(c Card) bool {
		return c.Rank == Two || c.Rank == Three
	}
	d := NewDeck(Filter(exclude))

	for _, c := range d {
		if exclude(c) {
			t.Errorf("Found an unexpected %s", c)
		}
	}
}

func TestCombine(t *testing.T) {
	d := NewDeck(Combine(NewDeck(), NewDeck()))
	if len(d) != 3*52 {
		t.Errorf("Combined deck of length %d (expected %d)", len(d), 3*52)
	}
}
