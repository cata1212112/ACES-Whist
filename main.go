package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

var Port = 5000

func main() {
	if err := openDatabase(); err != nil {
		log.Printf("Error opening database: %v", err)
	}
	defer closeDatabase() // close the database after main returns

	router := mux.NewRouter()
	router.Use(authMiddleware) // Adding the auth middleware to the router

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/login", loginGETHandler).Methods("GET")
	router.HandleFunc("/login", loginPOSTHandler).Methods("POST")
	router.HandleFunc("/setCookie", cookieTestHandler).Methods("GET")
	router.HandleFunc("/getCookies", getCookiesHandler).Methods("GET")
	router.HandleFunc("/register", registerPOSTHandler).Methods("POST")
	router.HandleFunc("/register", registerGETHandler).Methods("GET")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		renderError(w, http.StatusNotFound)
	})

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

	http.ListenAndServe(":" + fmt.Sprint(Port), router)
}
