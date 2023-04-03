package main

// Round
// playRound
// plays the round according to the rules

type Round struct {
	sumOfBids int
	trump     Card
	first     Card
}

func (round *Round) playRound(players *[]Player, deck *Deck, numberOfCards int) {
	round.sumOfBids = 0
	round.trump.value = -1
	round.first.value = -1

	for i := 0; i < 4; i++ {
		(*players)[i].giveCards(deck.giveCards(numberOfCards))
	}

	if numberOfCards < 8 {
		round.trump = deck.giveCards(1)[0]
	}
	sum := 0

	for i := 0; i < 3; i++ {
		(*players)[i].makeBid(false, 0, numberOfCards)
		(*players)[i].tricks = 0
		sum += (*players)[i].getBid()
	}
	(*players)[3].makeBid(true, sum, numberOfCards)
	(*players)[3].tricks = 0
	isFirst := 1
	for i := 0; i < numberOfCards; i++ {
		var winningCard Card
		var winningPlayer *Player
		winningPlayer = nil
		for i := 0; i < 4; i++ {
			var played Card
			if isFirst == 1 {
				if round.trump.value == -1 {
					played = (*players)[i].playCard(nil, nil)
				} else {
					played = (*players)[i].playCard(nil, &round.trump)
				}
				round.first = played
				winningCard = played
				winningPlayer = &(*players)[i]
				isFirst = 0
			} else {
				played = (*players)[i].playCard(&round.first, &round.trump)
			}
			var trumpCard *Card
			if round.trump.value == -1 {
				trumpCard = nil
			} else {
				trumpCard = &round.trump
			}
			if isFirst == 0 && played.compare(winningCard, trumpCard, round.first) {
				winningCard = played
				winningPlayer = &(*players)[i]
			}
		}
		winningPlayer.tricks++
	}

	for i := 0; i < 4; i++ {
		if (*players)[i].tricks == (*players)[i].bid {
			(*players)[i].score++
		}
	}
}
