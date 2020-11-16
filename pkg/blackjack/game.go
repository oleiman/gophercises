package blackjack

import (
	"fmt"
	"gophercises/pkg/deck"
)

type Game struct {
	Players []Player
	Deck    deck.Deck
}

type GameOption func(g *Game)

func Players(n int) GameOption {
	return func(g *Game) {
		for i := 0; i < n; i++ {
			g.Players = append(g.Players, NewHumanPlayer(i+1))
		}
	}
}

func NewGame(d deck.Deck, opts ...GameOption) *Game {
	g := Game{
		Players: make([]Player, 0, 2),
		Deck:    d,
	}
	for _, opt := range opts {
		opt(&g)
	}
	g.Players = append(g.Players, NewDealer())
	return &g
}

func (g *Game) Deal() {
	for _, p := range g.Players {
		var c deck.Card
		for i := 0; i < 2; i++ {
			c, g.Deck = g.Deck.Next()
			p.Accept(c)
		}
	}
}

func (g *Game) Dealer() *Dealer {
	return g.Players[len(g.Players)-1].(*Dealer)
}

func (g *Game) HumanPlayers() []Player {
	return g.Players[:len(g.Players)-1]
}

func (g *Game) Play() {
	for _, p := range g.Players {
		for p.Score() < 21 && p.TakeTurn() == Hit {
			var c deck.Card
			c, g.Deck = g.Deck.Next()
			p.Accept(c)
		}
		p.Stand()
		fmt.Println(p)
	}

	d := g.Dealer()
	for _, p := range g.HumanPlayers() {
		// TODO(oren): strict g.t.?
		p.(*Human).Won(d)
	}
}

func (g *Game) String() string {
	var str string
	for _, p := range g.Players {
		str = fmt.Sprintf("%s%s\n", str, p)
	}
	return str
}
