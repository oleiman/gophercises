package blackjack

import (
	"bufio"
	"fmt"
	"gophercises/pkg/deck"
	"os"
	"strings"
)

//go:generate stringer -type=Action

type Action int

const (
	Hit Action = iota
	Stand
)

type Player interface {
	Stand()
	Accept(c deck.Card)
	TakeTurn() Action
	Score() int
}

type Human struct {
	hand   Hand
	Winner bool
	ID     int
}

type Dealer struct {
	hand   Hand
	hidden Hand
}

func NewHumanPlayer(id int) *Human {
	return &Human{
		hand:   make(Hand, 0, 2),
		Winner: false,
		ID:     id,
	}
}

func (p *Human) Stand() {
	// do nothing
}

func (p *Human) Accept(c deck.Card) {
	p.hand = append(p.hand, c)
}

func (p *Human) TakeTurn() Action {
	fmt.Println(p)
	fmt.Printf("Player %d: Hit? ", p.ID)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if in[0] == 'n' {
		return Stand
	}
	return Hit
}

func (p *Human) Score() int {
	s, _ := p.hand.Score()
	return s
}

func (p *Human) Won(d *Dealer) bool {
	if p.Score() > 21 {
		return false
	} else if d.Score() > 21 || p.Score() == 21 || p.Score() > d.Score() {
		p.Winner = true
	}
	return p.Winner
}

func NewDealer() *Dealer {
	return &Dealer{
		hand:   make(Hand, 0, 2),
		hidden: make(Hand, 0, 1),
	}
}

func (p *Dealer) Stand() {
	p.hand = append(p.hidden, p.hand...)
	p.hidden = nil
}

func (p *Dealer) Accept(c deck.Card) {
	if len(p.hidden) == 0 {
		p.hidden = append(p.hidden, c)
	} else {
		p.hand = append(p.hand, c)
	}
}

func (p *Dealer) fullScore() (int, bool) {
	var tmp Hand = append(p.hand, p.hidden...)
	total, soft := tmp.Score()
	return total, soft
}

func (p *Dealer) TakeTurn() Action {
	total, soft17 := p.fullScore()
	if total <= 16 || soft17 {
		return Hit
	}
	return Stand
}

func (p *Dealer) Score() int {
	s, _ := p.hand.Score()
	return s
}

func (p Human) String() string {
	msg := ""
	if p.Score() == 21 && len(p.hand) == 2 {
		msg = " *BLACKJACK*"
	} else if p.Winner {
		msg = " *WINNER*"
	} else if p.Score() > 21 {
		msg = " *BUST*"
	}
	return fmt.Sprintf(
		"Human Player %d%s\n\tHand: %s\n\tScore: %d",
		p.ID, msg, p.hand, p.Score())
}

func (p Dealer) String() string {
	hidden := ""
	if len(p.hidden) > 0 {
		hidden = "HIDDEN, "
	}
	msg := ""
	if p.Score() == 21 && len(p.hand)+len(p.hidden) == 2 {
		msg = " *BLACKJACK*"
	} else if p.Score() > 21 {
		msg = " *BUST"
	}
	return fmt.Sprintf(
		"Dealer%s\n\tHand: %s%s\n\tScore: %d",
		msg, hidden, p.hand, p.Score())
}
