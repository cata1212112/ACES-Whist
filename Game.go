package main

import "fmt"

// Game
// func newGame() *Game
//		creates a new game
//	func (game *Game) addPlayer(player Player)
// 		used to add a player to the current game
// func (game *Game) play()
// 		plays the game according to the rules

type Game struct {
	numberOfCards []int
	deckOfCards   Deck
	players       []Player
	name          string
}

func newGame() *Game {
	game := new(Game)
	game.numberOfCards = []int{8, 8, 8, 8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8}
	game.deckOfCards = *NewDeck()
	return game
}

func (game *Game) addPlayer(player Player) {
	game.players = append(game.players, player)
}

func (game *Game) play() {
	for _, elem := range game.numberOfCards {
		for j := 0; j < 4; j++ {
			game.players[j].tricks = 0
		}
		round := new(Round)
		game.deckOfCards.index = 0
		game.deckOfCards.ShuffleDeck()
		round.playRound(&game.players, &game.deckOfCards, elem, game.name)
	}
	for j := 0; j < 4; j++ {
		fmt.Println(game.players[j].score)
	}
}
