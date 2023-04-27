package main

import (
	// "log"
	"fmt"
	"encoding/gob"
	"net/http"

	// "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte("tCP2QkKC2QO5NPukJLWbKfWzuaPgHcaNMPxfGC6bkj2U6KGrCN")) //super-secret-password :)

func cookieStoreInit() {
	store.Options.HttpOnly = true
	store.Options.Secure = true // requires secure HTTPS connection TODO: maybe set to false... IDK
	gob.Register(&User{})
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request URL and METHOD: %s	%s\n", r.URL.String(), r.Method)
        // Check if user is authenticated
        if !isAuthenticated(r) && r.URL.String() != "/" && r.URL.String() != "/login" && r.URL.String() != "/register"{
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func isAuthenticated(r *http.Request) bool {
    session, _ := store.Get(r, "session-id")
	_, ok := session.Values["username"]
	return ok
}

func checkLoginOK(userToLogin *UserLoginRequest) int {
	user, err := getUserByUsername(userToLogin.username)
	if err != nil {
		return http.StatusNotFound // There is no user with that username
	}

	if !checkPasswordHash(userToLogin.password, user.passwordHash) {
		return http.StatusUnauthorized // Password is incorrect
	}

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

func logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "session-id")
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
