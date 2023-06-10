package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	Name   string `json:"name"`
	bid    int
	Score  int `json:"score"`
	cards  []Card
	tricks int
}

func (player *Player) setName(name string) {
	player.Name = name
}

func (player *Player) getBid() int {
	return player.bid
}

func (player *Player) addScore(x int) {
	player.Score += x
}

func (player *Player) makeBid(isLast bool, sumBids int, numberOfCards int, gameID string) {
	/// asteapta un bid asteapta un post din forntend
	ok := 1
	for ok == 1 {
		// trimite pe canalul pe care jocul un masaj playerului cu numele name sa bage big

		command := map[string]interface{}{
			"method": "publish",
			"params": map[string]interface{}{
				"channel": gameID,
				"data": map[string]interface{}{
					"user": player.Name,
					"flag": "requestBid",
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

		fmt.Println("Please make bid " + player.Name + ":")
		gameMapMu.RLock()
		ch := gameMap[gameID]
		gameMapMu.RUnlock()

		// Wait for player input
		input := <-ch

		// Handle player input for the game
		fmt.Printf("Received input for game %d from player %d\n", input.GameID, input.Player.Name, input.Player.bid)
		// Update game state accordingly

		// sa dea publish pa canalul jocului cu noua stare

		if isLast == true && player.bid != numberOfCards-sumBids {
			ok = 0
		} else {
			ok = 0
		}
	}
}

func (player *Player) GiveCards(cards []Card) {
	player.cards = make([]Card, len(cards))
	copy(player.cards, cards)
}

func (player *Player) HasSuite(card Card) bool {
	for _, x := range player.cards {
		if x.Suite == card.Suite {
			return true
		}
	}
	return false
}

func (player *Player) GetValidCards(first *Card, trump *Card) []Card {
	validCards := make([]Card, len(player.cards))
	index := 0
	var hasFirst bool
	var hasTrump bool
	if first == nil {
		hasFirst = false
	} else {
		hasFirst = player.HasSuite(*first)
	}
	if trump == nil {
		hasTrump = false
	} else {
		hasTrump = player.HasSuite(*trump)
	}

	if hasFirst {
		for _, elem := range player.cards {
			if elem.Suite == first.Suite {
				validCards[index] = elem
				index++
			}
		}
	} else if hasTrump {

		for _, elem := range player.cards {
			if elem.Suite == trump.Suite {
				validCards[index] = elem
				index++
			}
		}
	} else {
		copy(validCards, player.cards)
	}
	cards := make([]Card, 0)
	for _, crd := range validCards {
		if crd.Value != 0 {
			cards = append(cards, crd)
		}
	}
	return cards
}

func (player *Player) playCard(first *Card, trump *Card, gameID string) Card {
	validCards := make([]Card, len(player.cards))
	validCards = player.GetValidCards(first, trump)
	fmt.Println("Cartile valide sunt")
	fmt.Println(validCards)

	var validCardsToSend PlayerCards
	validCardsToSend.Player = player.Name
	validCardsToSend.Cards = validCards
	jsonData, err := json.Marshal(validCardsToSend)
	os.Stdout.Write(jsonData)

	command := map[string]interface{}{
		"method": "publish",
		"params": map[string]interface{}{
			"channel": gameID,
			"data": map[string]interface{}{
				"data": jsonData,
				"flag": "validCards",
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

	gameMapMu.RLock()
	ch := gameMap[gameID]
	gameMapMu.RUnlock()

	input := <-ch
	var j int
	for i, c := range player.cards {
		if c.Suite == input.playedCard.Suite && c.Value == input.playedCard.Value {
			j = i
		}
	}
	if j >= 0 && j < len(player.cards) {
		player.cards = append(player.cards[:j], player.cards[j+1:]...)
	}

	command1 := map[string]interface{}{
		"method": "publish",
		"params": map[string]interface{}{
			"channel": gameID,
			"data": map[string]interface{}{
				"who":   player.Name,
				"value": input.playedCard.Value,
				"suite": input.playedCard.Suite,
				"flag":  "playedCard",
			},
		},
	}

	dataA1, err := json.Marshal(command1)
	if err != nil {
		panic(err)
	}
	req1, err := http.NewRequest("POST", "http://localhost:8000/api", bytes.NewBuffer(dataA1))
	if err != nil {
		panic(err)
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce")
	client1 := http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	return input.playedCard
}
