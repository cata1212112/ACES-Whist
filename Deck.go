package main

import (
	"math/rand"
	"sort"
)

// Deck
// holds a deck of cards
// func newDeck() *Deck
//		creates a new deck
// func (deck *Deck) shuffleDeck()
//		shuffles the deck of card
// func (deck *Deck) giveCards(i int) []Card
//		used to remove i cards from the deck

const J = 12
const Q = 13
const K = 14
const A = 15

const deckSize = 32

type Deck struct {
	cards []Card
	index int
}

func NewDeck() *Deck {
	deck := new(Deck)
	deck.index = 0
	for i := 7; i < 11; i++ {
		for _, s := range []suites{HEARTS, SPADES, DIAMONDS, CLUBS} {
			deck.cards = append(deck.cards, Card{value: i, suite: s})
		}
	}

	for _, s := range []suites{HEARTS, SPADES, DIAMONDS, CLUBS} {
		deck.cards = append(deck.cards, Card{value: J, suite: s})
	}
	for _, s := range []suites{HEARTS, SPADES, DIAMONDS, CLUBS} {
		deck.cards = append(deck.cards, Card{value: Q, suite: s})
	}
	for _, s := range []suites{HEARTS, SPADES, DIAMONDS, CLUBS} {
		deck.cards = append(deck.cards, Card{value: K, suite: s})
	}
	for _, s := range []suites{HEARTS, SPADES, DIAMONDS, CLUBS} {
		deck.cards = append(deck.cards, Card{value: A, suite: s})
	}
	return deck
}

func (deck *Deck) ShuffleDeck() {
	sort.Slice(deck.cards, func(i, j int) bool {
		return rand.Intn(1000) < rand.Intn(1000)
	})
}

func (deck *Deck) GiveCards(i int) []Card {
	/// nu cred ca vom avea nevoie de mai multe carti cat sunt in pachet decat daca e un bug

	if i > deckSize-deck.index+1 {
		return nil
	}
	aux := make([]Card, i)
	copy(aux, deck.cards[deck.index:deck.index+i])
	deck.index += i
	return aux
}

func (deck *Deck) Equals(otherDeck Deck) bool {
	for i, c := range deck.cards {
		if c.suite != otherDeck.cards[i].suite && c.value != otherDeck.cards[i].value {
			return false
		}
	}
	return true
}
