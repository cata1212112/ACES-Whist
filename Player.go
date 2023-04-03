package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//Player
// func (player *Player) setName(name string)
// 		name setter
// func (player *Player) getBid() int
// 		bid getter
// func (player *Player) addScore(x int)
// 		adds x score to a player
// func (player *Player) makeBid(isLast bool, sumBids int, numberOfCards int)
// 		requests the player to make a bid
// func (player *Player) giveCards(cards []Card)
//		used to five cards to the player
//	func (player *Player) hasSuite(card Card) bool
//		checks if the player has a card with the same suite as card
// func (player *Player) playCard(first *Card, trump *Card) Card
//		requests the player to play a card

type Player struct {
	name   string
	bid    int
	score  int
	cards  []Card
	tricks int
}

func (player *Player) setName(name string) {
	player.name = name
}

func (player *Player) getBid() int {
	return player.bid
}

func (player *Player) addScore(x int) {
	player.score += x
}

func (player *Player) makeBid(isLast bool, sumBids int, numberOfCards int) {
	/// asteapta un bid asteapta un post din forntend
	ok := 1
	for ok == 1 {
		fmt.Println("Please make bid " + player.name + ":")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		player.bid, _ = strconv.Atoi(text)

		if isLast == true && player.bid != numberOfCards-sumBids {
			ok = 0
		} else {
			ok = 0
		}
	}
}

func (player *Player) giveCards(cards []Card) {
	player.cards = make([]Card, len(cards))
	copy(player.cards, cards)
}

func (player *Player) hasSuite(card Card) bool {
	for _, x := range player.cards {
		if x.suite == card.suite {
			return true
		}
	}
	return false
}

func (player *Player) playCard(first *Card, trump *Card) Card {

	validCards := make([]Card, len(player.cards))
	index := 0
	var hasFirst bool
	var hasTrump bool
	if first == nil {
		hasFirst = false
	} else {
		hasFirst = player.hasSuite(*first)
	}
	if trump == nil {
		hasTrump = false
	} else {
		hasTrump = player.hasSuite(*trump)
	}

	if hasFirst {
		for _, elem := range player.cards {
			if elem.suite == first.suite {
				validCards[index] = elem
				index++
			}
		}
	} else if hasTrump {

		for _, elem := range player.cards {
			if elem.suite == trump.suite {
				validCards[index] = elem
				index++
			}
		}
	} else {
		copy(validCards, player.cards)
	}

	ok := 1
	var i int
	fmt.Println("Cartile valide sunt")
	fmt.Println(validCards)
	for ok == 1 {
		fmt.Println("Please choose card from the valid ones " + player.name + ":")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		i, _ = strconv.Atoi(text)
		if i >= 0 && i <= len(validCards) {
			ok = 0
		}
	}

	return validCards[i]
}
