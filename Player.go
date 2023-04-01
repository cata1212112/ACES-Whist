package ACES

import "fmt"

type Player struct {
	name  string
	bid   int
	score int
	cards []Card
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
		fmt.Scan("%d", &player.bid)
		if isLast && player.bid != numberOfCards-sumBids {
			ok = 0
		}
	}
}

func (player *Player) giveCards(cards []Card) {
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

func (player *Player) playCard(first Card, trump Card) Card {

	var validCards []Card
	hasFirst := player.hasSuite(first)
	hasTrump := player.hasSuite(trump)

	if hasFirst {
		for _, elem := range player.cards {
			if elem.suite == first.suite {
				validCards = append(validCards, elem)
			}
		}
	} else if hasTrump {
		for _, elem := range player.cards {
			if elem.suite == trump.suite {
				validCards = append(validCards, elem)
			}
		}
	} else {
		copy(validCards, player.cards)
	}

	ok := 1
	var i int
	for ok == 1 {
		fmt.Println("Please choose card from the valid ones" + player.name + ":")
		fmt.Scan("%d", &i)
		if i >= 0 && i <= len(validCards) {
			ok = 0
		}
	}

	return validCards[i]
}
