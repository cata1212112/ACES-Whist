package main

import (
	"log"
	// "fmt"
	"encoding/gob"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte("tCP2QkKC2QO5NPukJLWbKfWzuaPgHcaNMPxfGC6bkj2U6KGrCN")) //super-secret-password :)

func cookieStoreInit() {
	store.Options.HttpOnly = true
	store.Options.Secure = true // requires secure HTTPS connection TODO: maybe set to false... IDK
	gob.Register(&User{})
}

func login(userToLogin *UserLoginRequest) int {
	user, err := getUserByUsername(userToLogin.username)
	if err != nil {
		return http.StatusNotFound // There is no user with that username
	}

	if !checkPasswordHash(userToLogin.password, user.passwordHash) {
		return http.StatusUnauthorized // Password is incorrect
	}

	// TODO: do the login process with the user and set the session_id field to something

	return http.StatusOK // replace with cookie set and session initialization
}

func register(userToCreate *UserCreate) int {
	if _, err := getUserByUsername(userToCreate.username); err == nil {
		return http.StatusConflict // There is already an user with that username
	}

	if _, err := getUserByEmail(userToCreate.email); err == nil {
		return http.StatusConflict // There is already an user with that email
	}

	if userToCreate.password != userToCreate.confirmPassword {
		return http.StatusUnauthorized // Passwords do not match
	}

	if passwordHash, err := hashPassword(userToCreate.password); err != nil {
		return http.StatusInternalServerError // Error during password hashing  <- maybe delete this if/else statement
	} else {
		userToCreate.password = passwordHash
	}

	if err := createUser(userToCreate); err != nil {
		return http.StatusInternalServerError // Database error when creating user
	}

	return http.StatusOK // Success
}

func logout() {
	// remove the cookie and destroy the session
}

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}

func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read basic auth information
		usr, _, ok := r.BasicAuth()

		// if there is no basic auth (no matter which credentials)
		if !ok {
			errMsg := "Authentication error!"
			// return a 403 forbidden
			http.Error(w, errMsg, http.StatusForbidden)
			log.Println(errMsg)

			// stop processing route
			return
		}

		// let's assume we check the user against a database to get
		// his admin-right and put this to the request context
		context.Set(r, "isAdmin", true)

		// else continue processing
		log.Printf("User %s logged in.", usr)
		next(w, r)
	}
}
