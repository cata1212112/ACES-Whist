package main

func main() {
	var player1 Player
	var player2 Player
	var player3 Player
	var player4 Player
	game := newGame()
	player1.setName("1")
	player2.setName("2")
	player3.setName("3")
	player4.setName("4")

	game.addPlayer(player1)
	game.addPlayer(player2)
	game.addPlayer(player3)
	game.addPlayer(player4)

	game.play()
}
