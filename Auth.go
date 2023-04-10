package main

import (
	"log"
    "fmt"
    _ "github.com/lib/pq"
)

func login() {
   
}

func register(userToCreate *UserCreate) bool {
	if _, err := getUserByUsername(userToCreate.username); err == nil {
		fmt.Printf("There is already a user with username %v", userToCreate.username)
		return false 		// replace with status code
	}

	if passwordHash, err := hashPassword(userToCreate.password); err != nil {
		fmt.Printf("Error hashing password: %v", err)
		return false	// replace with status code or delete case
	} else {
		userToCreate.password = passwordHash
	}

	if err := createUser(userToCreate); err != nil {
		log.Printf("Error creating user: %v", err)
		return false 		// replace with status code
	}
	
	return true			// replace with status code
}

func logout() {

}