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
	home, err := os.ReadFile("collector-page/login/login.html")
	if err != nil {
		//TODO add an unavailable page
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("{\"msg\": \"page unavailable at the moment\"}"))
	}
	w.Write(home)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var post LoginPost
	fmt.Println(r.Body)

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	user, err := IsCorrect(post.Username, post.Password)
	fmt.Println(post)

	//return response
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "msg": "User Unauthorized"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("{\"code\": 202, \"authLvl\": " + strconv.Itoa(user.Level) + " \"}"))
	fmt.Println(user)
}

func init() {
	ConnectToDatabaseUsers()
}

func main() {
	r := mux.NewRouter()

	//user login area
	r.HandleFunc(usersLogin.String(), HomeHandler).Methods("GET")
	r.HandleFunc(usersLogin.String(), LoginHandler).Methods("POST")

	http.ListenAndServe(":8080", r)
}
