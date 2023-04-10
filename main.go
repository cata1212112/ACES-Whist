package main

import (
	"log"
)

func main() {
	if err := openDatabase(); err != nil {
		log.Printf("Error opening database: %v", err)
	}
	defer closeDatabase()		// close the database after main returns

	// var player1 Player
	// var player2 Player
	// var player3 Player
	// var player4 Player
	// game := newGame()
	// player1.setName("1")
	// player2.setName("2")
	// player3.setName("3")
	// player4.setName("4")

	// game.addPlayer(player1)
	// game.addPlayer(player2)
	// game.addPlayer(player3)
	// game.addPlayer(player4)

	// game.play()
	testUser := UserCreate{
		email: "email@gmail.com",
		username: "username", 
		password: "password",
	}

	if err := register(&testUser); err == false {
		log.Printf("Couldn't create user: %v", testUser.username)
	} else {
		log.Printf("Created user: %v", testUser.username)
	}

}
