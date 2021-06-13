package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type LoginPost struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("collector-page/login.html")
	if err != nil {
		//TODO add an unavailable page
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write("{\"msg\": \"page unavailable at the moment\"")
	}
	w.Write(home)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var post LoginPost
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)
	_ = json.NewDecoder(r.Body).Decode(&post)

	user, err := IsCorrect(post.Username, post.Password)
	fmt.Println(post)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "msg": "User Unauthorized"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("{\"code\": 202, \"authLvl\": " + strconv.Itoa(user.Level) + " \""))
	fmt.Println(user)
}

func init() {
	ConnectToDatabaseUsers()
}

func main() {
	r := mux.NewRouter()
	//login area
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/", LoginHandler).Methods("POST")

	//collectors endpoints
	
	http.ListenAndServe(":8080", r)
}
