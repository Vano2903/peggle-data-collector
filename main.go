package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func GetPfp(w http.ResponseWriter, r *http.Request) {
	var post Post
	fmt.Println(r.Body)

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	user, err := QueryUser(post.Username, post.Password)
	fmt.Println(post)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}
	imgJson := fmt.Sprintf(`{"url":"%s"}`, user.PfpUrl)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(imgJson))
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

//TODO let user insert a date like 2021-2-4 and convert it to 2021-02-04 otherwhise it will return and error and it can be annoying
//TODO documentation cause error messages are becoming a bit too long XD
func SeachGameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ciaoo")
	r.ParseForm()
	sortByDate := true
	fmt.Println(r.Form)
	lim := 25
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
		case "upload": //the url query will be < >before the date and it will search dates before, after or equal to the date setted
			for _, ds := range v {
				dateElements := strings.Split(ds, "-")
				if len(dateElements) != 3 {
					PrintErr(w, "date not defined correctly")
					return
				}
				dateString := ds[1:] + "T00:00:00+00:00"
				d, err := time.Parse(time.RFC3339, dateString)
				if err != nil {
					PrintErr(w, err.Error())
				}
				switch ds[0] {
				case '>':
					q := bson.D{{"$match", bson.D{{"videoData.uploadDate", bson.M{"$gt": primitive.NewDateTimeFromTime(d)}}}}}
					queries = append(queries, q)
				case '<':
					q := bson.D{{"$match", bson.D{{"videoData.uploadDate", bson.M{"$lt": primitive.NewDateTimeFromTime(d)}}}}}
					queries = append(queries, q)
				default:
					PrintErr(w, "date search operand not correct, must use '<', '>'")
					return
				}
			}
		case "points": //the url will be s,r. to define if search in all the games made by syn or red and so, ro to search in the overall stats
			for _, ps := range v {
				pointsElements := strings.Split(ps, "-")
				if len(pointsElements) != 3 && len(pointsElements) != 2 {
					PrintErr(w, "points has an incorrect format")
					return
				}

				switch pointsElements[0] {
				case "r":
					points, err := strconv.Atoi(pointsElements[2])
					if err != nil {
						PrintErr(w, err.Error())
						return
					}
					switch pointsElements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.points", bson.M{"$gt": points}}}, {{"stats.redez.g2.points", bson.M{"$gt": points}}}, {{"stats.redez.g3.points", bson.M{"$gt": points}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.points", bson.M{"$lt": points}}}, {{"stats.redez.g2.points", bson.M{"$lt": points}}}, {{"stats.redez.g3.points", bson.M{"$lt": points}}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "incorrect operator in points search")
						return
					}
				case "s":
					points, err := strconv.Atoi(pointsElements[2])
					if err != nil {
						PrintErr(w, err.Error())
						return
					}
					switch pointsElements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.points", bson.M{"$gt": points}}}, {{"stats.synergo.g2.points", bson.M{"$gt": points}}}, {{"stats.synergo.g3.points", bson.M{"$gt": points}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.points", bson.M{"$lt": points}}}, {{"stats.synergo.g2.points", bson.M{"$lt": points}}}, {{"stats.synergo.g3.points", bson.M{"$lt": points}}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "incorrect operator in points search")
						return
					}
				case "ro":
					switch pointsElements[1] {
					case ">":
						points, err := strconv.Atoi(pointsElements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.redez.overall.tPoints", bson.M{"$gt": points}}}}}
						queries = append(queries, q)
					case "<":
						points, err := strconv.Atoi(pointsElements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.redez.overall.tPoints", bson.M{"$lt": points}}}}}
						queries = append(queries, q)
					case "max":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.redez.overall.tPoints": -1}}}
						queries = append(queries, q)
						sortByDate = false
					case "min":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.redez.overall.tPoints": 1}}}
						queries = append(queries, q)
						sortByDate = false
					default:
						PrintErr(w, "incorrect operator in points search")
						return
					}
				case "so":
					switch pointsElements[1] {
					case ">":
						points, err := strconv.Atoi(pointsElements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.synergo.overall.tPoints", bson.M{"$gt": points}}}}}
						queries = append(queries, q)
					case "<":
						points, err := strconv.Atoi(pointsElements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.synergo.overall.tPoints", bson.M{"$lt": points}}}}}
						queries = append(queries, q)
					case "max":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.synergo.overall.tPoints": -1}}}
						queries = append(queries, q)
						sortByDate = false
					case "min":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.synergo.overall.tPoints": 1}}}
						queries = append(queries, q)
						sortByDate = false
					default:
						PrintErr(w, "incorrect operator in points search")
						return
					}
				default:
					PrintErr(w, "points search operand not correct")
					return
				}
			}
		case "n-25":
			for _, ps := range v {
				n25Elements := strings.Split(ps, "-")
				if len(n25Elements) != 3 && len(n25Elements) != 2 {
					PrintErr(w, "n-25 has an incorrect format")
					return
				}
				fmt.Println(n25Elements)
				switch n25Elements[0] {
				case "r":
					points, err := strconv.Atoi(n25Elements[2])
					if err != nil {
						PrintErr(w, err.Error())
						return
					}
					switch n25Elements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.n-25", bson.M{"$gt": points}}}, {{"stats.redez.g2.n-25", bson.M{"$gt": points}}}, {{"stats.redez.g3.n-25", bson.M{"$gt": points}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.n-25", bson.M{"$lt": points}}}, {{"stats.redez.g2.n-25", bson.M{"$lt": points}}}, {{"stats.redez.g3.n-25", bson.M{"$lt": points}}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "incorrect operator in n-25 search")
						return
					}
				case "s":
					points, err := strconv.Atoi(n25Elements[2])
					if err != nil {
						PrintErr(w, err.Error())
						return
					}
					switch n25Elements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.n-25", bson.M{"$gt": points}}}, {{"stats.synergo.g2.n-25", bson.M{"$gt": points}}}, {{"stats.synergo.g3.n-25", bson.M{"$gt": points}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.n-25", bson.M{"$lt": points}}}, {{"stats.synergo.g2.n-25", bson.M{"$lt": points}}}, {{"stats.synergo.g3.n-25", bson.M{"$lt": points}}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "incorrect operator in n-25 search")
						return
					}
				case "ro":
					switch n25Elements[1] {
					case ">":
						points, err := strconv.Atoi(n25Elements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.redez.overall.t-25", bson.M{"$gt": points}}}}}
						queries = append(queries, q)
					case "<":
						points, err := strconv.Atoi(n25Elements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.redez.overall.t-25", bson.M{"$lt": points}}}}}
						queries = append(queries, q)
					case "max":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.redez.overall.t-25": -1}}}
						queries = append(queries, q)
						sortByDate = false
					case "min":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.redez.overall.t-25": 1}}}
						queries = append(queries, q)
						sortByDate = false
					default:
						PrintErr(w, "incorrect operator in n-25 search")
						return
					}
				case "so":
					switch n25Elements[1] {
					case ">":
						points, err := strconv.Atoi(n25Elements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.synergo.overall.t-25", bson.M{"$gt": points}}}}}
						queries = append(queries, q)
					case "<":
						points, err := strconv.Atoi(n25Elements[2])
						if err != nil {
							PrintErr(w, err.Error())
							return
						}
						q := bson.D{{"$match", bson.D{{"stats.synergo.overall.t-25", bson.M{"$lt": points}}}}}
						queries = append(queries, q)
					case "max":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.synergo.overall.t-25": -1}}}
						queries = append(queries, q)
						sortByDate = false
					case "min":
						lim = 1
						q := bson.D{{"$sort", bson.M{"stats.synergo.overall.t-25": 1}}}
						queries = append(queries, q)
						sortByDate = false
					default:
						PrintErr(w, "incorrect operator in n-25 search")
						return
					}
				default:
					PrintErr(w, "n-25 search operand not correct")
					return
				}
			}
		case "val-fe":
			for _, fes := range v {
				exfeElements := strings.Split(fes, "-")
				if len(exfeElements) != 3 {
					PrintErr(w, "extreme fever not defined correctly")
					return
				}
				if exfeElements[2] != "0" && exfeElements[2] != "5000" && exfeElements[2] != "25000" && exfeElements[2] != "50000" {
					PrintErr(w, "extreme fever's value is invalid, must use 0, 5000, 25000, 50000 only")
					return
				}
				valFe, err := strconv.Atoi(exfeElements[2])
				if err != nil {
					PrintErr(w, err.Error())
					return
				}

				switch exfeElements[0] {
				case "r":
					switch exfeElements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.valFE", bson.M{"$gt": valFe}}}, {{"stats.redez.g2.valFE", bson.M{"$gt": valFe}}}, {{"stats.redez.g3.valFE", bson.M{"$gt": valFe}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.valFE", bson.M{"$lt": valFe}}}, {{"stats.redez.g2.valFE", bson.M{"$lt": valFe}}}, {{"stats.redez.g3.valFE", bson.M{"$lt": valFe}}}}}}}}
						queries = append(queries, q)
					case "e":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.redez.g1.valFE", valFe}}, {{"stats.redez.g2.valFE", valFe}}, {{"stats.redez.g3.valFE", valFe}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "extreme fever search operand not correct, must use '<', '>', 'e'")
						return
					}
				case "s":
					switch exfeElements[1] {
					case ">":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.valFE", bson.M{"$gt": valFe}}}, {{"stats.synergo.g2.valFE", bson.M{"$gt": valFe}}}, {{"stats.synergo.g3.valFE", bson.M{"$gt": valFe}}}}}}}}
						queries = append(queries, q)
					case "<":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.valFE", bson.M{"$lt": valFe}}}, {{"stats.synergo.g2.valFE", bson.M{"$lt": valFe}}}, {{"stats.synergo.g3.valFE", bson.M{"$lt": valFe}}}}}}}}
						queries = append(queries, q)
					case "e":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{{{"stats.synergo.g1.valFE", valFe}}, {{"stats.synergo.g2.valFE", valFe}}, {{"stats.synergo.g3.valFE", valFe}}}}}}}
						queries = append(queries, q)
					default:
						PrintErr(w, "extreme fever search operand not correct, must use '<', '>', 'e'")
						return
					}
				}
			}
		case "character":
			//TODO
		case "limit":
			if len(v) != 1 {
				PrintErr(w, "can't define more than a limit option")
				return
			}
			val, err := strconv.Atoi(v[0])
			if err != nil {
				PrintErr(w, err.Error())
				return
			}
			if val <= 0 {
				PrintErr(w, "definte a positive number grater than 0 for the limit")
				return
			}
			lim = val
		default:
			PrintErr(w, "invalid parameter")
			return
		}
	}
	if sortByDate {
		q := bson.D{{"$sort", bson.M{"videoData.uploadDate": 1}}}
		queries = append(queries, q)
	}
	limit := bson.D{{"$limit", lim}}
	queries = append(queries, limit)

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

	//get url for user's pfp
	r.HandleFunc(usersPfp.String(), GetPfp).Methods("POST")

	//commit area
	r.HandleFunc(getCommits.String(), CommitHandler).Methods("POST")

	//game area
	r.HandleFunc(games.String(), SeachGameHandler).Methods("GET")

	log.Println("starting on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
