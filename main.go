package main

import (
	"encoding/json"
	"fmt"
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
	home, err := os.ReadFile("pages/login/login.html")
	if err != nil {
		//TODO add an unavailable page
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("{\"msg\": \"page unavailable at the moment\"}"))
		return
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
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "msg": "User Unauthorized"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	// w.Write([]byte("{\"code\": 202, \"authLvl\": " + strconv.Itoa(user.Level) + " \"}"))

	var page []byte
	switch user.Level {
	case 0:
		page, err = os.ReadFile("pages/lvl0/index.html")
	case 1:
		page, err = os.ReadFile("pages/lvl1/index.html")
	case 2:
		page, err = os.ReadFile("pages/lvl2/index.html")
	}
	if err != nil {
		//TODO add an unavailable page
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("{\"msg\": \"page unavailable at the moment\"}"))
		return
	}
	w.Write(page)
	fmt.Println("the user: ", user.User, " just logged in, the auth lvl is: ", user.Level)
}

func init() {
	ConnectToDatabaseUsers()
}

func main() {
	r := mux.NewRouter()

	//user login area
	r.HandleFunc(usersLogin.String(), HomeHandler).Methods("GET")
	r.HandleFunc(usersLogin.String(), LoginHandler).Methods("POST")
	log.Println("starting on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
