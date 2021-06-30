package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Username string `json: "username, omitempty"`
	Password string `json: "password, omitempty"`
	Year     int    `json: "year, omitempty"`
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
	var post Post
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

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	//read user credentials
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	//get param from url
	param := mux.Vars(r)["param"]
	fmt.Println("param ", param)

	switch param {
	//return the ammount of total commits made by the user
	case "totCommits":
		tot, err := GetTotalCommits(post.Username, post.Password)
		if err != nil {
			PrintErr(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"totalCommits": %d}`, tot)))

	//return all the years a user has made at least 1 commit
	case "years":
		years, err := GetCommitsYear(post.Username, post.Password)
		fmt.Println(years)
		if err != nil {
			PrintErr(w, err.Error())
			return
		}
		var resp string
		for _, y := range years {
			resp += strconv.Itoa(y) + ";"
		}
		if len(resp) > 0 {
			resp = resp[:len(resp)-1]
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(resp))
		return

	//return all the commits made by the user over the selected year
	case "year":
		commits, err := GetCommitsByYear(post.Username, post.Password, post.Year)
		if err != nil {
			PrintErr(w, err.Error())
			return
		}
		j, err := json.Marshal(commits)
		if err != nil {
			PrintInternalErr(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return

	//increment the ammount of commits by 1 on the user's document
	case "add":
		err := AddCommit(post.Username, post.Password)
		if err != nil {
			PrintErr(w, err.Error())
			return
		}
		w.Write([]byte("commit registered"))
		return

	//return bad request
	default:
		PrintErr(w, "invalid parameter")
		return
	}
}

func SeachGameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ciaoo")
	r.ParseForm()

	fmt.Println(r.Form)

	var queries []bson.D
	// var result []User
	for k, v := range r.Form {
		switch k {
		case "id":
			q := bson.D{{"$match", bson.D{{"videoData.id", bson.M{"$in": v}}}}}
			queries = append(queries, q)
		case "title":
			if len(v) != 1 {
				PrintErr(w, "you can't query over multiple titles (yet)")
				return
			}
			q := bson.D{{"$match", bson.D{{"videoData.title", bson.D{{"$regex", primitive.Regex{Pattern: v[0], Options: "i"}}}}}}}
			queries = append(queries, q)
		case "wonBy":
			vi, err := ConvertToSliceInt(v)
			if err != nil {
				PrintErr(w, err.Error())
				return
			}
			q := bson.D{{"$match", bson.D{{"wonBy", bson.M{"$in": vi}}}}}
			queries = append(queries, q)
		case "upload":
			//TODO
		case "points":

		case "n-25":

		case "valEF":

		case "character":

		}
	}

	fmt.Println(queries)
	fmt.Println(QueryGames(queries))
}

func main() {
	r := mux.NewRouter()
	//statics
	r.PathPrefix(statics.String()).Handler(http.StripPrefix(statics.String(), http.FileServer(http.Dir("static/"))))

	//user login area
	r.HandleFunc(usersLogin.String(), HomeHandler).Methods("GET")
	r.HandleFunc(usersLogin.String(), LoginHandler).Methods("POST")

	//commit area
	r.HandleFunc(getCommits.String(), CommitHandler).Methods("POST")

	//game area
	r.HandleFunc(games.String(), SeachGameHandler).Methods("GET")

	log.Println("starting on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
