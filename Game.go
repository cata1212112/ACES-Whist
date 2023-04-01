package ACES

type Game struct {
	numberOfCards []int
	deckOfCards   Deck
	players       []Player
}

func newGame() *Game {
	game := new(Game)
	game.numberOfCards = []int{8, 8, 8, 8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8}
	game.deckOfCards = *newDeck()
	return game
}

func (game *Game) addPlayer(player Player) {
	game.players = append(game.players, player)
}

func (game *Game) play() {

}
