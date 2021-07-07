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
	Username string `json:"username, omitempty"`
	Password string `json:"password, omitempty"`
	Year     int    `json:"year, omitempty"`
	Id       string `json:"id, omitempty"`
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

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	user, err := QueryUser(post.Username, post.Password)
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

	//read post body
	_ = json.NewDecoder(r.Body).Decode(&post)

	//check if user is correct
	user, err := IsCorrect(post.Username, post.Password)

	//return response
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
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
	log.Println("the user:", user.User, " just logged in, the auth lvl is:", user.Level)
}

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	//read user credentials
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	//get param from url
	param := mux.Vars(r)["param"]

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
	//parse url query
	r.ParseForm()
	sortByDate := true
	lim := 25
	var queries []bson.D
	//all options :D
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
			for _, chs := range v {
				char := chs[1:]
				switch chs[0] {
				case 's':
					switch char {
					case "cas":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "castoro"}},
							{{"stats.synergo.g2.valFE", "castoro"}},
							{{"stats.synergo.g3.valFE", "castoro"}}}}}}}
						queries = append(queries, q)
					case "uni":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "unicorno"}},
							{{"stats.synergo.g2.valFE", "unicorno"}},
							{{"stats.synergo.g3.valFE", "unicorno"}}}}}}}
						queries = append(queries, q)
					case "zuc":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "zucca"}},
							{{"stats.synergo.g2.valFE", "zucca"}},
							{{"stats.synergo.g3.valFE", "zucca"}}}}}}}
						queries = append(queries, q)
					case "gat":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "gatto"}},
							{{"stats.synergo.g2.valFE", "gatto"}},
							{{"stats.synergo.g3.valFE", "gatto"}}}}}}}
						queries = append(queries, q)
					case "ali":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "alieno"}},
							{{"stats.synergo.g2.valFE", "alieno"}},
							{{"stats.synergo.g3.valFE", "alieno"}}}}}}}
						queries = append(queries, q)
					case "gra":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "granchio"}},
							{{"stats.synergo.g2.valFE", "granchio"}},
							{{"stats.synergo.g3.valFE", "granchio"}}}}}}}
						queries = append(queries, q)
					case "gir":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "girasole"}},
							{{"stats.redez.g2.valFE", "girasole"}},
							{{"stats.redez.g3.valFE", "girasole"}}}}}}}
						queries = append(queries, q)
					case "dra":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "drago"}},
							{{"stats.synergo.g2.valFE", "drago"}},
							{{"stats.synergo.g3.valFE", "drago"}}}}}}}
						queries = append(queries, q)
					case "con":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "coniglio"}},
							{{"stats.synergo.g2.valFE", "coniglio"}},
							{{"stats.synergo.g3.valFE", "coniglio"}}}}}}}
						queries = append(queries, q)
					case "guf":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "gufo"}},
							{{"stats.synergo.g2.valFE", "gufo"}},
							{{"stats.synergo.g3.valFE", "gufo"}}}}}}}
						queries = append(queries, q)
					case "sep":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.synergo.g1.character", "seppia"}},
							{{"stats.synergo.g2.valFE", "seppia"}},
							{{"stats.synergo.g3.valFE", "seppia"}}}}}}}
						queries = append(queries, q)
					}
				case 'r':
					switch char {
					case "cas":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "castoro"}},
							{{"stats.redez.g2.valFE", "castoro"}},
							{{"stats.redez.g3.valFE", "castoro"}}}}}}}
						queries = append(queries, q)
					case "uni":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "unicorno"}},
							{{"stats.redez.g2.valFE", "unicorno"}},
							{{"stats.redez.g3.valFE", "unicorno"}}}}}}}
						queries = append(queries, q)
					case "zuc":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "zucca"}},
							{{"stats.redez.g2.valFE", "zucca"}},
							{{"stats.redez.g3.valFE", "zucca"}}}}}}}
						queries = append(queries, q)
					case "gat":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "gatto"}},
							{{"stats.redez.g2.valFE", "gatto"}},
							{{"stats.redez.g3.valFE", "gatto"}}}}}}}
						queries = append(queries, q)
					case "ali":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "alieno"}},
							{{"stats.redez.g2.valFE", "alieno"}},
							{{"stats.redez.g3.valFE", "alieno"}}}}}}}
						queries = append(queries, q)
					case "gra":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "granchio"}},
							{{"stats.redez.g2.valFE", "granchio"}},
							{{"stats.redez.g3.valFE", "granchio"}}}}}}}
						queries = append(queries, q)
					case "gir":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "girasole"}},
							{{"stats.redez.g2.valFE", "girasole"}},
							{{"stats.redez.g3.valFE", "girasole"}}}}}}}
						queries = append(queries, q)
					case "dra":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "drago"}},
							{{"stats.redez.g2.valFE", "drago"}},
							{{"stats.redez.g3.valFE", "drago"}}}}}}}
						queries = append(queries, q)
					case "con":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "coniglio"}},
							{{"stats.redez.g2.valFE", "coniglio"}},
							{{"stats.redez.g3.valFE", "coniglio"}}}}}}}
						queries = append(queries, q)
					case "guf":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "gufo"}},
							{{"stats.redez.g2.valFE", "gufo"}},
							{{"stats.redez.g3.valFE", "gufo"}}}}}}}
						queries = append(queries, q)
					case "sep":
						q := bson.D{{"$match", bson.D{{"$or", []bson.D{
							{{"stats.redez.g1.character", "seppia"}},
							{{"stats.redez.g2.valFE", "seppia"}},
							{{"stats.redez.g3.valFE", "seppia"}}}}}}}
						queries = append(queries, q)
					}
				}
			}
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

	results, err := QueryGames(queries)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}
	j, err := json.Marshal(results)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

//given a youtube url will search in database and if it's not found will return a message saying the video is not yet registered in the database
//otherwise will return a json of the data of the game found
func CheckGameHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := []bson.D{bson.D{{"$match", bson.D{{"videoData.id", params["id"]}}}}}
	results, err := QueryGames(q)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	//if len is 0 then nothing was found and let the user know that this id is not in the database
	if len(results) == 0 {
		j := `{"msg": "the id is not in the database"}`
		w.Write([]byte(j))
		return
	}

	//return the stats of the game if found (need it for the update of a game)
	j, err := json.Marshal(results)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}

	w.Write(j)
}

func AddGameHandler(w http.ResponseWriter, r *http.Request) {
	//read user credentials
	var post Game
	json.NewDecoder(r.Body).Decode(&post)
	//TODO function to check if some values are missing in the post request
	err := post.VD.GetYoutubeDataFromId(post.VD.Id)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}
	fmt.Println(post)
	id, err := AddGame(post)
	if err != nil {
		PrintErr(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"msg":"Game added correctly and it's id is: %s"}`, id)))
}

func UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	//updated game in post
	var post Game
	json.NewDecoder(r.Body).Decode(&post)
	//id of the game to update
	idGameUpdate := mux.Vars(r)["id"]
	found, _ := CheckIfExist(idGameUpdate)
	if !found {
		PrintErr(w, "game not found")
		return
	}
	err := FullUpdateGame(idGameUpdate, post)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"game updated correctly"}`))
}

func DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	err := DeleteGame(post.Id)
	if err != nil {
		PrintInternalErr(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"game deleted correctly"}`))
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
	r.HandleFunc(checkGame.String(), CheckGameHandler).Methods("GET")
	r.HandleFunc(addGame.String(), AddGameHandler).Methods("POST")
	r.HandleFunc(updateGame.String(), UpdateGameHandler).Methods("POST")
	r.HandleFunc(deleteGame.String(), DeleteGameHandler).Methods("POST")

	log.Println("starting on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
