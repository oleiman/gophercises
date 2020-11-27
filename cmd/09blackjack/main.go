package main

import (
	"fmt"
	"gophercises/pkg/blackjack"
	"gophercises/pkg/deck"
)

func main() {
	d := deck.NewDeck(deck.Shuffle)
	g := blackjack.NewGame(d, blackjack.Players(2))
	g.Deal()
	fmt.Println(g.Dealer())
	g.Play()
	fmt.Println("\nFINAL: ")
	fmt.Print(g)
}
