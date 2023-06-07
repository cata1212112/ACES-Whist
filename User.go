package main

import (
	"errors"
	"log"
	"net/http"
)

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	passwordHash string
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreate struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func createUser(user *UserCreate) error {
	if err := DB.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password).Err(); err != nil {
		return err
	}
	return nil
}

func getUserByUsername(username string) (*User, error) {
	var myUser User
	if err := DB.QueryRow("SELECT username, email, password FROM users WHERE username = $1", username).Scan(&myUser.Username, &myUser.Email, &myUser.passwordHash); err != nil {
		return nil, err
	}
	return &myUser, nil
}

func getUserByEmail(email string) (*User, error) {
	var myUser User
	if err := DB.QueryRow("SELECT username, email, password FROM users WHERE email = $1", email).Scan(&myUser.Username, &myUser.Email, &myUser.passwordHash); err != nil {
		return nil, err
	}
	return &myUser, nil
}

func getFriendsOfUser(user User) ([]User, error) {
	var user_id int
	if err := DB.QueryRow("select id from users where username = $1", user.Username).Scan(&user_id); err != nil {
		return []User{}, err
	}
	var friends []User
	rows, err := DB.Query("select username from users where id in ((select user_id1 from are_friends where user_id2 = $1 and confirmed_1 and confirmed_2) union (select user_id2 from are_friends where user_id1 = $1 and confirmed_1 and confirmed_2))", user_id)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return []User{}, err
		}
		friends = append(friends, User{Username: username})
	}

	return friends, nil
}

func getFriendRequestsOfUser(user User) ([]User, error) {
	var user_id int
	if err := DB.QueryRow("select id from users where username = $1", user.Username).Scan(&user_id); err != nil {
		return []User{}, err
	}
	var friend_requests []User
	rows, err := DB.Query("select username from users where id in ((select user_id1 from are_friends where user_id2 = $1 and confirmed_1 and not confirmed_2) union (select user_id2 from are_friends where user_id1 = $1 and not confirmed_1 and confirmed_2))", user_id)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return []User{}, err
		}
		friend_requests = append(friend_requests, User{Username: username})
	}
	return friend_requests, nil
}

func getUsersNotRelatedToMe(user User) ([]User, error) {
	var user_id int
	if err := DB.QueryRow("select id from users where username = $1", user.Username).Scan(&user_id); err != nil {
		return []User{}, err
	}

	var usersNotRelated []User
	rows, err := DB.Query("select username from users where id not in (select user_id1 from are_friends where user_id2 = $1 union select user_id2 from are_friends where user_id1 = $1)", user_id)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return []User{}, err
		}
		usersNotRelated = append(usersNotRelated, User{Username: username})
	}
	return usersNotRelated, nil
}

func areFriends(user1 User, user2 User) (bool, error) {
	var result bool
	var err error
	var user1_id int
	var user2_id int
	if err = DB.QueryRow("select id from users where username = $1", user1.Username).Scan(&user1_id); err != nil {
		return false, err
	}

	if err = DB.QueryRow("select id from users where username = $1", user2.Username).Scan(&user2_id); err != nil {
		return false, err
	}

	if err = DB.QueryRow("select exists(select 1 from are_friends where ((user_id1, user_id2) = ($1, $2) or (user_id1, user_id2) = ($2, $1)) and confirmed_1 and confirmed_2)", user1_id, user2_id).Scan(&result); err != nil {
		return false, err
	}

	return result, nil
}

func sendFriendRequest(sender User, receiver User) error {
	var sender_id int
	var receiver_id int

	if err := DB.QueryRow("select id from users where username = $1", sender.Username).Scan(&sender_id); err != nil {
		return err
	}

	if err := DB.QueryRow("select id from users where username = $1", receiver.Username).Scan(&receiver_id); err != nil {
		return err
	}

	if sender_id == receiver_id {
		return errors.New("Cannot send friend request to yourself!")
	}

	if _, err := DB.Exec("insert into are_friends(user_id1, user_id2, confirmed_1, confirmed_2) values($1, $2, true, false)", sender_id, receiver_id); err != nil {
		return err
	}

	return nil
}

func acceptFriendRequest(accepter User, other User) error {
	var accepter_id int
	var other_id int

	if err := DB.QueryRow("select id from users where username = $1", accepter.Username).Scan(&accepter_id); err != nil {
		return err
	}

	if err := DB.QueryRow("select id from users where username = $1", other.Username).Scan(&other_id); err != nil {
		return err
	}

	if accepter_id == other_id {
		return errors.New("Cannot accept friend request from yourself!")
	}

	result, err := DB.Exec("update are_friends set confirmed_2 = true where user_id1 = $1 and user_id2 = $2", other_id, accepter_id)
	if err != nil {
		return err
	}
	affected_rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected_rows == 0 {
		return errors.New("No friend request to accept!")
	}
	return nil
}

func rejectFriendRequest(rejecter User, other User) error {
	var rejecter_id int
	var other_id int

	if err := DB.QueryRow("select id from users where username = $1", rejecter.Username).Scan(&rejecter_id); err != nil {
		return err
	}

	if err := DB.QueryRow("select id from users where username = $1", other.Username).Scan(&other_id); err != nil {
		return err
	}

	if rejecter_id == other_id {
		return errors.New("Cannot reject a friend request from yourself!")
	}

	result, err := DB.Exec("delete from are_friends where (user_id1, userid_2) = ($1, $2) or (user_id1, userid_2) = ($2, $1)", other_id, rejecter_id)
	if err != nil {
		return err
	}

	affected_rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected_rows == 0 {
		return errors.New("No friend request to reject!")
	}
	return nil
}

func unfriend(user1 User, user2 User) error {
	are_friends, err := areFriends(user1, user2)
	if err != nil {
		return err
	}
	if !are_friends {
		return errors.New("Users are not friends!")
	}

	var user1_id int
	var user2_id int

	if err := DB.QueryRow("select id from users where username = $1", user1.Username).Scan(&user1_id); err != nil {
		return err
	}

	if err := DB.QueryRow("select id from users where username = $1", user2.Username).Scan(&user2_id); err != nil {
		return err
	}

	if user1_id == user2_id {
		return errors.New("Cannot unfriend yourself!")
	}

	if _, err := DB.Exec("delete from are_friends where (user_id1, user_id2) = ($1, $2) or (user_id1, user_id2) = ($2, $1)", user1_id, user2_id); err != nil {
		return err
	}
	return nil
}

// Test function for user creation
func testUserCreate() { // modify to test the create of other users
	testUser := UserCreate{
		Email:           "email@gmail.com",
		Username:        "username",
		Password:        "password",
		ConfirmPassword: "password",
	}

	if status := register(&testUser); status != http.StatusOK {
		log.Printf("Couldn't create user: %v (Status code: %v)\n", testUser.Username, status)
	} else {
		log.Printf("Created user: %v\n", testUser.Username)
	}
}
