package main

import (
	// "fmt"
	"log"

	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	if err := openDatabase(); err != nil {
		log.Printf("Error opening database: %v", err)
	}
	defer closeDatabase()		// close the database after main returns

	router := mux.NewRouter()

	router.HandleFunc("/", testHandler).Methods("GET")
	router.HandleFunc("/login", loginGETHandler).Methods("GET")
	router.HandleFunc("/login", loginPOSTHandler).Methods("POST")

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



	// Login Tests
	// userToLogin := UserLoginRequest {
	// 	username: "username",
	// 	password: "password",
	// }

	// fmt.Println("First login (OK)")
	// if login(&userToLogin) {
	// 	fmt.Println("Login successful")
	// } else {
	// 	fmt.Println("Login failed")
	// }

	// fmt.Println()
	// fmt.Println("Second login (user does not exist)")
	// tempUsername := userToLogin.username
	// userToLogin.username = "badUsername"
	// if login(&userToLogin) {
	// 	fmt.Println("Login successful")
	// } else {
	// 	fmt.Println("Login failed")
	// }

	// userToLogin.username = tempUsername
	// userToLogin.password = "badPassword"
	// fmt.Println()
	// fmt.Println("Third login (incorrect password)")
	// if login(&userToLogin) {
	// 	fmt.Println("Login successful")
	// } else {
	// 	fmt.Println("Login failed")
	// }

	http.ListenAndServe(":5000", router)
}
