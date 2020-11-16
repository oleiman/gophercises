package blackjack

import (
	"fmt"
	"gophercises/pkg/deck"
	"testing"
)

func TestScoreHand(t *testing.T) {
	hand := Hand{
		deck.Card{Suit: deck.Hearts, Rank: deck.Two},
		deck.Card{Suit: deck.Hearts, Rank: deck.Four},
		deck.Card{Suit: deck.Hearts, Rank: deck.King},
	}

	s, _ := hand.Score()
	if s != 16 {
		t.Errorf("Incorrect score: %d (expected %d)", s, 16)
	}

	hand = append(hand, deck.Card{Suit: deck.Hearts, Rank: deck.Ace})
	s, soft := hand.Score()
	if s != 17 || soft {
		t.Errorf("Incorrect score: %d (soft: %t) (expected hard 17)", s, soft)
	}

	hand = Hand{
		deck.Card{Suit: deck.Clubs, Rank: deck.Six},
		deck.Card{Suit: deck.Clubs, Rank: deck.Ace},
	}

	s, soft = hand.Score()
	if s != 17 || !soft {
		t.Errorf("Incorrect score: %d (soft: %t) (expected soft 17)", s, soft)
	}
}

func TestHandString(t *testing.T) {
	hand := Hand{
		deck.Card{Suit: deck.Clubs, Rank: deck.Jack},
		deck.Card{Suit: deck.Clubs, Rank: deck.Ace},
	}
	expected := "Jack of Clubs, Ace of Clubs"
	if hand.String() != expected {
		t.Errorf(`"%s" should be "%s"`, hand, expected)
	}
}

func TestHumanPlayer(t *testing.T) {
	cards := []deck.Card{
		{Suit: deck.Hearts, Rank: deck.Two},
		{Suit: deck.Hearts, Rank: deck.Four},
	}
	pID := 1
	p := NewHumanPlayer(pID)
	for _, c := range cards {
		p.Accept(c)
	}

	if p.ID != pID {
		t.Errorf("Player ID: %d (expected %d)", p.ID, pID)
	}

	if len(p.hand) != 2 {
		t.Errorf("Human player has %d cards (expect %d)", len(p.hand), 2)
	}

	strRep := "Human Player 1\n\tHand: Two of Hearts, Four of Hearts\n\tScore: 6"
	if p.String() != strRep {
		t.Errorf(`"%s" should be "%s"`, p, strRep)
	}
}

func TestDealer(t *testing.T) {
	cards := []deck.Card{
		{Suit: deck.Hearts, Rank: deck.Two},
		{Suit: deck.Hearts, Rank: deck.Four},
	}
	p := NewDealer()
	for _, c := range cards {
		p.Accept(c)
	}

	if len(p.hand) != 1 {
		t.Errorf("Human player has %d cards (expect %d)", len(p.hand), 1)
	}

	strRep := "Dealer\n\tHand: HIDDEN, Four of Hearts\n\tScore: 4"
	if p.String() != strRep {
		t.Errorf(`"%s" should be "%s"`, p, strRep)
	}

	p.Stand()

	if len(p.hand) != 2 {
		t.Errorf("Human player has %d cards (expect %d)", len(p.hand), 1)
	}

	strRep = "Dealer: Two of Hearts, Four of Hearts"
	strRep = "Dealer\n\tHand: Two of Hearts, Four of Hearts\n\tScore: 6"
	if p.String() != strRep {
		t.Errorf(`"%s" should be "%s"`, p, strRep)
	}
}

func TestDealNewGame(t *testing.T) {
	n := 1
	d := deck.NewDeck()
	g := NewGame(d, Players(n))

	if len(g.Players)-1 != n {
		t.Errorf("Found %d players (expected %d)", len(g.Players)-1, n)
	}
}

func TestDealGame(t *testing.T) {
	dck := deck.NewDeck()
	_, dck = dck.Next()
	g := NewGame(dck, Players(1))
	var p *Human = g.Players[0].(*Human)
	var d *Dealer = g.Players[1].(*Dealer)

	// discard the first ace for easier accounting
	g.Deal()

	pScore := int(deck.Two) + int(deck.Three)
	if p.Score() != pScore {
		t.Errorf("Player score is %d (expected %d)", p.Score(), pScore)
	}

	dScore := int(deck.Five)
	if d.Score() != dScore {
		fmt.Println(d)
		t.Errorf("Dealer score is %d (expected %d)", d.Score(), dScore)
	}

	// dealer reveals other card, increasing visible score
	d.Stand()
	dScore += int(deck.Four)
	if d.Score() != dScore {
		fmt.Println(d)
		t.Errorf("Dealer score is %d (expected %d)", d.Score(), dScore)
	}
}
