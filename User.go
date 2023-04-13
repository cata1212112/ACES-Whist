package main

import (
	"log"
	"net/http"
)

type User struct {
	username string `json:"username"`
	email    string `json:"email"`
	passwordHash string 
}

type UserLoginRequest struct {
	username string `json:"username"`
	password string `json:"password"`
	session_id string `json:"session-id"`
}

type UserCreate struct {
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
	confirmPassword string `json:"confirmPassword"`
}

func createUser(user *UserCreate) error {
	if err := DB.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.username, user.email, user.password).Err(); err != nil {
		return err
	}
	return nil
}

func getUserByUsername(username string) (*User, error) {
	var myUser User
	if err := DB.QueryRow("SELECT username, email, password FROM users WHERE username = $1", username).Scan(&myUser.username, &myUser.email, &myUser.passwordHash); err != nil {
		return nil, err;
	}
	return &myUser, nil;
}

func getUserByEmail(email string) (*User, error) {
	var myUser User
	if err := DB.QueryRow("SELECT username, email, password FROM users WHERE email = $1", email).Scan(&myUser.username, &myUser.email, &myUser.passwordHash); err!= nil {
        return nil, err;
    }
	return &myUser, nil;
}

func testUserCreate() {		// modify to test the create of other users
	testUser := UserCreate{
		email: "email@gmail.com",
		username: "username", 
		password: "password",
		confirmPassword: "password",
	}

	if status := register(&testUser); status != http.StatusOK {
		log.Printf("Couldn't create user: %v (Status code: %v)\n", testUser.username, status)
	} else {
		log.Printf("Created user: %v\n", testUser.username)
	}
}