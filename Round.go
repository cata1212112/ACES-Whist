package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Round
// playRound
// plays the round according to the rules

type Round struct {
	sumOfBids int
	trump     Card
	first     Card
}

func (round *Round) playRound(players *[]Player, deck *Deck, numberOfCards int, gameID string) {
	fmt.Println("intrat")
	round.sumOfBids = 0
	round.trump.Value = -1
	round.first.Value = -1

	for i := 0; i < 4; i++ {
		(*players)[i].GiveCards(deck.GiveCards(numberOfCards))
	}

	if numberOfCards < 8 {
		round.trump = deck.GiveCards(1)[0]
	}

	/// trimit trump + cartile fiecauria
	var gameDTO GameDTO
	gameDTO.Trump = round.trump
	for i := 0; i < 4; i++ {
		gameDTO.Players[i].Player = (*players)[i].Name
		gameDTO.Players[i].Cards = (*players)[i].cards
	}
	jsonData, err := json.Marshal(gameDTO)
	os.Stdout.Write(jsonData)

	command := map[string]interface{}{
		"method": "publish",
		"params": map[string]interface{}{
			"channel": gameID,
			"data": map[string]interface{}{
				"data": jsonData,
				"flag": "carti_joc",
			},
		},
	}

	dataA, err := json.Marshal(command)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8000/api", bytes.NewBuffer(dataA))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	sum := 0

	for i := 0; i < 3; i++ {
		fmt.Println("cer bid playerului ", i)
		(*players)[i].makeBid(false, 0, numberOfCards, gameID)

		(*players)[i].tricks = 0
		sum += (*players)[i].getBid()
	}
	(*players)[3].makeBid(true, sum, numberOfCards, gameID)
	(*players)[3].tricks = 0
	for i := 0; i < numberOfCards; i++ {
		var winningCard Card
		var winningPlayer *Player
		winningPlayer = nil
		isFirst := 1
		for i := 0; i < 4; i++ {
			var played Card
			if isFirst == 1 {
				if round.trump.Value == -1 {
					played = (*players)[i].playCard(nil, nil, gameID)
				} else {
					played = (*players)[i].playCard(nil, &round.trump, gameID)
				}
				round.first = played
				winningCard = played
				winningPlayer = &(*players)[i]
				isFirst = 0
			} else {
				played = (*players)[i].playCard(&round.first, &round.trump, gameID)
			}
			var trumpCard *Card
			if round.trump.Value == -1 {
				trumpCard = nil
			} else {
				trumpCard = &round.trump
			}
			if isFirst == 0 && played.Compare(winningCard, trumpCard, round.first) {
				winningCard = played
				winningPlayer = &(*players)[i]
			}
		}

		command := map[string]interface{}{
			"method": "publish",
			"params": map[string]interface{}{
				"channel": gameID,
				"data": map[string]interface{}{
					"flag": "endHand",
				},
			},
		}

		dataA, err := json.Marshal(command)
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest("POST", "http://localhost:8000/api", bytes.NewBuffer(dataA))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		winningPlayer.tricks++
	}

	for i := 0; i < 4; i++ {
		if (*players)[i].tricks == (*players)[i].bid {
			(*players)[i].Score += 10 + (*players)[i].bid
		} else {
			//(*players)[i].Score -= math.Abs((*players)[i].tricks - (*players)[i].bid)
		}
	}
}
