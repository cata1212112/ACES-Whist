package main

import(
	"net/http"
	"html/template"
	"log"
	"fmt"
	"github.com/gorilla/sessions"
)


type PageData struct {
	Username interface{}
}

type ErrorPageData struct {
	Error interface{}
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	if _, ok := session.Values["username"]; !ok {		// user not logged in
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
		username: r.MultipartForm.Value["username"][0],
		email: r.MultipartForm.Value["email"][0],
		password: r.MultipartForm.Value["password"][0],
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
	
	r.ParseMultipartForm(32 << 20)		// this gets the content of the form in form-data format
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
		Path: "/",
        MaxAge:   3600 * 24 * 7,	// the cookie will expire in a week
        HttpOnly: true,
	}

    session.Values["username"] = userToLogin.username
    err = session.Save(r, w)
	if err != nil {
		renderError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loginGETHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	if _, ok := session.Values["username"]; ok {		// if the user is already logged in, then don't let them see the login form
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
		Path: "/",
		MaxAge: 5,		// test here
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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound)
}