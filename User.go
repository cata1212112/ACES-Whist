package main

type User struct {
	id       int64  `json:"id"`
	username string `json:"username"`
	email    string `json:"email"`
}

type UserCreate struct {
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
}

type UserRequest struct {
	username string `json:"username"`
	password string `json:"password"`
	token    string `json:"token"`
}

func createUser(user *UserCreate) error {
	if err := DB.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.username, user.email, user.password).Err(); err != nil {
		return err
	}
	return nil
}

func getUserByUsername(username string) (*User, error) {
	var myUser User
	if err := DB.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&myUser); err != nil {
		return nil, err;
	}
	return &myUser, nil;
}

// func getAllusers() ([]User, error) {
// 	var (
// 		rows *sql.Rows
// 	 	err error
// 		result []User
// 	)

// 	if rows, err = DB.Query("SELECT * FROM users"); err != nil {
// 		return nil, err
// 	}

// 	result = make([]User, len(rows))

// }
