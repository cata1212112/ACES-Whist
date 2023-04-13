package main

import(
	"net/http"
	"html/template"
	"log"
	"fmt"
)


func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head></head><body>AM REUSIT SA FAC CEVA!</body></html>"))
}

func loginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()  -->  r.Form    <=  if you want to get the login form in x-xxx-form-urlencoded format
	
	r.ParseMultipartForm(32 << 20)		// this gets the content of the form in form-data format
	fmt.Println("Form: ", r.MultipartForm)
	var userToLogin UserLoginRequest = UserLoginRequest{
		username: r.MultipartForm.Value["username"][0],
	    password: r.MultipartForm.Value["password"][0],
		session_id: "",
    }


	t, _ := template.ParseFiles("error-template.html")
	if err := t.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}

	if status := login(&userToLogin); status != http.StatusOK {
		r.Response.StatusCode = status
		r.Response.Header.Add("session_id", userToLogin.session_id)		// TODO: dintre astea doua
		w.Header().Add("session_id", userToLogin.session_id)			// cred ca trebuie doar una pastrata
		return
	}
	
}

func loginGETHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("login-template.html")
	if err := t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}