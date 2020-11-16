package blackjack

import (
	"fmt"
	"gophercises/pkg/deck"
)

type Hand []deck.Card

func (h Hand) Score() (int, bool) {
	score := 0
	nAces := 0
	soft := false
	for _, c := range h {
		if c.Rank == deck.Ace {
			nAces++
			score += int(c.Rank)
		} else if c.Rank == deck.Jack || c.Rank == deck.Queen || c.Rank == deck.King {
			score += 10
		} else {
			score += int(c.Rank)
		}
	}
	for nAces > 0 && score+10 <= 21 {
		soft = true
		score += 10
		nAces--
	}
	return score, soft
}

func (h Hand) String() string {
	var str string
	for _, c := range h {
		if len(str) == 0 {
			str = fmt.Sprintf("%s", c)
		} else {
			str = fmt.Sprintf("%s, %s", str, c)
		}
	}
	return str
}
