package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
)

type PageData struct {
	Username interface{}
	Token    interface{}
}

type LobbyData struct {
	Username  interface{}
	Token     interface{}
	LobbyName interface{}
}

type ErrorPageData struct {
	Error interface{}
}

var connectedClients = make(map[string]bool)
var lobbies = make(map[string][]string)
var mutex = &sync.Mutex{}

func isTokenAlreadyConnected(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return connectedClients[token]
}

func addConnectedClient(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	connectedClients[token] = true
}

func removeConnectedClient(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(connectedClients, token)
}

func renderError(w http.ResponseWriter, err interface{}) {
	t, _ := template.ParseFiles("error-template.html")
	data := ErrorPageData{Error: err}
	t.Execute(w, data)
}

// Just a simple test to see how these Handlers work
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head></head><body>AM REUSIT SA FAC CEVA!</body></html>"))
}

func lobbiesHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("lobbies.html")
	session, _ := store.Get(r, "session-id")
	if isTokenAlreadyConnected(session.Values["token"].(string)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pageData := PageData{Username: session.Values["username"], Token: session.Values["token"]}
	t.Execute(w, pageData)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	if _, ok := session.Values["username"]; !ok { // user not logged in
		t, _ := template.ParseFiles("login-template.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("profile.html")
		pageData := PageData{Username: session.Values["username"]}
		t.Execute(w, pageData)
	}
}

func registerPOSTHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var userToCreate UserCreate = UserCreate{
		username:        r.MultipartForm.Value["username"][0],
		email:           r.MultipartForm.Value["email"][0],
		password:        r.MultipartForm.Value["password"][0],
		confirmPassword: r.MultipartForm.Value["confirm-password"][0],
	}

	if err := register(&userToCreate); err != http.StatusOK {
		renderError(w, err)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func registerGETHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("register-template.html")
	t.Execute(w, nil)
}

func loginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()  -->  r.Form    <=  if you want to get the login form in x-www-form-urlencoded format

	r.ParseMultipartForm(32 << 20) // this gets the content of the form in form-data format
	fmt.Println("Form: ", r.MultipartForm)
	var userToLogin UserLoginRequest = UserLoginRequest{
		username: r.MultipartForm.Value["username"][0],
		password: r.MultipartForm.Value["password"][0],
	}

	if status := checkLoginOK(&userToLogin); status != http.StatusOK {
		renderError(w, status)
		return
	}

	session, err := store.Get(r, "session-id")
	if err != nil {
		renderError(w, err)
		return
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // the cookie will expire in a week
		HttpOnly: true,
	}

	session.Values["username"] = userToLogin.username
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userToLogin.username})
	tokenString, errToken := token.SignedString([]byte("da4a14bb-f4d7-4a32-90b3-15fb080d3937"))
	if errToken != nil {
		fmt.Println("Erorr generating token", errToken.Error())
		return
	}
	session.Values["token"] = tokenString
	err = session.Save(r, w)
	if err != nil {
		renderError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func onConnect(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-id")
	if err != nil {
		return
	}
	tokenString := session.Values["token"].(string)
	println(tokenString)
	if isTokenAlreadyConnected(tokenString) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	addConnectedClient(tokenString)
}

func onDisconnect(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString := session.Values["token"].(string)
	println(tokenString)
	if isTokenAlreadyConnected(tokenString) {
		removeConnectedClient(tokenString)
	}
	w.WriteHeader(http.StatusOK)
}

func loginGETHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	if _, ok := session.Values["username"]; ok { // if the user is already logged in, then don't let them see the login form
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("login-template.html")
	if err := t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := logout(w, r); err != nil {
		fmt.Fprintf(w, "Error during logout: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func cookieTestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 5, // test here
		// MaxAge:   3600 * 24 * 7,	// the cookie will expire in a week
		HttpOnly: true,
	}

	session.Values["foo"] = "bar"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCookiesHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		fmt.Fprintf(w, "Cookie: %v=%v\n", cookie.Name, cookie.Value)
	}
}

func getAllLobbies(w http.ResponseWriter, r *http.Request) {

	//'Authorization': 'apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce',
	println("ok")
	command := map[string]interface{}{
		"method": "channels",
		"params": map[string]interface{}{},
	}
	data, err := json.Marshal(command)
	if err != nil {
		println(err.Error())
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8000/api", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error(), resp)
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var result map[string]interface{}
	err = decoder.Decode(&result)
	if err != nil {
		println(err.Error())
		return
	}
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		println(err.Error())
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
	println(string(jsonResult))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound)
}

func lobbyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lobbyName := params["lobbyName"]
	t, _ := template.ParseFiles("lobby.html")
	session, _ := store.Get(r, "session-id")
	pageData := LobbyData{Username: session.Values["username"], Token: session.Values["token"], LobbyName: lobbyName}
	t.Execute(w, pageData)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	var data struct {
		User    string `json:"user" :"user"`
		Message string `json:"message" :"message"`
		Channel string `json:"channel" :"channel"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}
	command := map[string]interface{}{
		"method": "publish",
		"params": map[string]interface{}{
			"channel": data.Channel,
			"data": map[string]interface{}{
				"user":    data.User,
				"message": data.Message,
			},
		},
	}

	dataA, err := json.Marshal(command)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8000/api", bytes.NewBuffer(dataA))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey a3d9c270-52df-45f8-9a66-a1bb8e9e04ce")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func getStyleFile(w http.ResponseWriter, r *http.Request) {
	println("mergi")
	vars := mux.Vars(r)
	filename := vars["filename"]
	println(filename)
	filePath := "style/" + filename
	http.ServeFile(w, r, filePath)
}

func getCard(w http.ResponseWriter, r *http.Request) {
	println("mergi")
	vars := mux.Vars(r)
	filename := vars["filename"]
	println(filename)
	filePath := "deckOfCards/SVG-cards-1.3/" + filename
	http.ServeFile(w, r, filePath)
}

func addToLobbyHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Lobby string `json:"lobby"`
		Name  string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	lobby := requestData.Lobby
	if _, ok := lobbies[lobby]; !ok {
		lobbies[lobby] = make([]string, 0)
	}
	lobbies[lobby] = append(lobbies[lobby], requestData.Name)
}

func lobbyMembers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lobby := params["lobbyName"]
	println(lobbies)
	members, ok := lobbies[lobby]
	if !ok {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}
	responseData, err := json.Marshal(members)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	println(members)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func removeFromLobbyHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Lobby string `json:"lobby"`
		Name  string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	lobby := requestData.Lobby
	if _, ok := lobbies[lobby]; !ok {
		lobbies[lobby] = make([]string, 0)
	}
	for i := 0; i < len(lobbies[lobby]); i++ {
		if lobbies[lobby][i] == requestData.Name {
			lobbies[lobby] = append(lobbies[lobby][:i], lobbies[lobby][i+1:]...)
			break
		}
	}
}
