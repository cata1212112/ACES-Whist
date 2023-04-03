package main

import "fmt"

// Card
// holds the value and suite of a card
// func (card *Card) compare(winningCard Card, trump *Card, first Card) bool
// 		used to compare the current card with the winning card -> return true if the current card is better
//															   -> return false, otherwise

// func (card *Card) equals(otherCard Card) bool
//		used to compare the suites of two cards

type suites int

const (
	HEARTS   = 1
	CLUBS    = 2
	DIAMONDS = 3
	SPADES   = 4
)

type Card struct {
	suite suites
	value int
}

func (card *Card) compare(winningCard Card, trump *Card, first Card) bool {
	if trump != nil && card.suite != trump.suite && card.suite != first.suite {
		return false
	}
	if card.suite == winningCard.suite {
		return card.value > winningCard.value
	}
	if card.suite == first.suite && winningCard.suite == trump.suite {
		return false
	}
	if winningCard.suite == first.suite && card.suite == trump.suite {
		return true
	}
	return true
}

func (card *Card) equals(otherCard Card) bool {
	return card.suite == otherCard.suite
}

func (card *Card) showCard() { // nu avem nevoie de functia asta
	fmt.Printf("%d %d\n", card.suite, card.value)
}
