package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type LoginPost struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("collector-page/login.html")
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(home)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var post LoginPost
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(post)
	if IsAuthorised(post.Username, post.Password) {

	}
}

func main() {
	r := mux.NewRouter()
	//login area
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/", LoginHandler).Methods("POST")

	//admin endpoints

	// http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", r)
}
