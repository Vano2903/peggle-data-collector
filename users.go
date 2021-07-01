package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientUser     *mongo.Client
	ctxUser        context.Context
	collectionUser *mongo.Collection
)

// type DateU struct {
// 	Day   int `bson:"day, omitempty" json: "day, omitempty"`
// 	Month int `bson:"month, omitempty" json: "month, omitempty"`
// 	Year  int `bson:"year, omitempty" json: "year, omitempty"`
// }

// func (d *DateU) Today() {
// 	t := time.Now()
// 	y, m, da := t.Date()
// 	d.Year = y
// 	d.Month = int(m)
// 	d.Day = da
// }

type Commit struct {
	Totals    int                `bson:"totals, omitempty" json: "totals, omitempty"`
	CreatedAt primitive.DateTime `bson:"date, omitempty" json: "date, omitempty"`
}

type Stats struct {
	TotalCommits int      `bson:"totals, omitempty" json: "totals, omitempty"`
	Commits      []Commit `bson:"commits, omitempty" json: "commits, omitempty"`
}

//TODO pfpUrl and gitusername can be used as search parameters
type User struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json: "id, omitempty"`
	User        string             `bson:"user, omitempty" json: "user, omitempty"`
	GitUserName string             `bson:"git_user_name, omitempty" json: "git_user_name, omitempty"`
	PfpUrl      string             `bson:"pfp_url, omitempty" json: "pfp_url, omitempty"`
	Pass        string             `bson:"password, omitempty"  json: "password, omitempty"`
	Level       int                `bson:"authLevel, omitempty" json: "authLevel, omitempty"` //0 admin (every priviliges), 1 every power relative to the games, 2 just adding
	Stats       Stats              `bson:"stats, omitempty" json:"stats, omitempty"`
}

func (x User) IsStructureEmpty() bool {
	return reflect.DeepEqual(x, User{})
}

//will connect to database on user's collectionn
func ConnectToDatabaseUsers() error {
	//get context
	ctxUser, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.9:27017")
	clientUser, err := mongo.Connect(ctxUser, clientOptions)
	if err != nil {
		return err
	}

	//check if connection is established
	err = clientUser.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	//assign to the global variable "collection" the users' collection
	collectionUser = clientUser.Database("qdss-peggle").Collection("users")
	return nil
}

//return all the years in which something has been commited
func GetCommitsYear(user, pass string) ([]int, error) {
	u, err := QueryUser(user, pass)
	fmt.Println(u)
	if err != nil {
		return nil, err
	}

	var years []int
	for _, c := range u.Stats.Commits {
		fmt.Println(years)
		date := c.CreatedAt.Time()
		if !Contains(years, date.Year()) {
			years = append(years, date.Year())
		}
	}
	fmt.Println(years)
	return years, nil
}

//return all commits made by an user
func GetCommits(user, pass string) ([]Commit, error) {
	u, err := QueryUser(user, pass)
	if err != nil {
		return nil, err
	}
	return u.Stats.Commits, nil
}

func GetTotalCommits(user, pass string) (int, error) {
	u, err := QueryUser(user, pass)
	if err != nil {
		return -1, err
	}
	return u.Stats.TotalCommits, nil
}

//return all the commits made by an user in one year
func GetCommitsByYear(user, pass string, year int) ([]Commit, error) {
	u, err := QueryUser(user, pass)
	if err != nil {
		return nil, err
	}

	var commitsFoundByYear []Commit
	for _, c := range u.Stats.Commits {
		y := c.CreatedAt.Time().Year()
		if y == year {
			commitsFoundByYear = append(commitsFoundByYear, c)
		}
	}
	return commitsFoundByYear, nil
}

//increment commits total nuber and increment totals commit of the day or create a new day if it's the first one of the day
func AddCommit(user, pass string) error {
	User, err := QueryUser(user, pass)
	if err != nil {
		return err
	}
	type date struct {
		d int
		m time.Month
		y int
	}
	today := time.Now()
	t := date{today.Day(), today.Month(), today.Year()}
	// today := today.Date()
	for ind, com := range User.Stats.Commits {
		supp := com.CreatedAt.Time()
		d := date{supp.Day(), supp.Month(), supp.Year()}
		if reflect.DeepEqual(d, t) {
			User.Stats.TotalCommits += 1
			User.Stats.Commits[ind].Totals += 1
			update := bson.M{"stats": User.Stats}
			fmt.Println(update)
			UpdateUser(user, pass, update)
			return nil
		}
	}

	com := Commit{1, primitive.NewDateTimeFromTime(today)}
	User.Stats.Commits = append(User.Stats.Commits, com)

	User.Stats.TotalCommits += 1
	update := bson.M{"stats": User.Stats}
	fmt.Println(update)
	UpdateUser(user, pass, update)
	return nil
}

//check if user exist in database and will return empty struct if not found, on the other hand will return the User informations
//? I belive there is a bettere way to do this but rn i dont really know
func IsCorrect(user, pass string) (User, error) {
	//search in database
	cur, err := collectionUser.Find(ctxUser, bson.M{"user": user, "password": pass})
	if err != nil {
		return User{}, err
	}
	defer cur.Close(ctxUser)
	var userFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &userFound); err != nil {
		return User{}, err
	}

	//check if user exist
	if len(userFound) != 0 {
		if userFound[0].User == user && userFound[0].Pass == pass {
			return userFound[0], nil
		}
		return User{}, errors.New("incorrect credentials")
	} else {
		return User{}, errors.New("no user found")
	}
}

//return the user based on username and password
func QueryUser(user, pass string) (User, error) {
	query := bson.M{"user": user, "password": pass}
	cur, err := collectionUser.Find(ctxUser, query)
	if err != nil {
		return User{}, err
	}
	defer cur.Close(ctxUser)
	var userFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &userFound); err != nil {
		return User{}, err
	}
	return userFound[0], nil
}

//return a []User given a bson.M query (possible queries: user, authLevel)
func QueryUsers(query bson.M) ([]User, error) {
	if _, ok := query["password"]; ok {
		return nil, errors.New("you can't query over passwords")
	}
	cur, err := collectionUser.Find(ctxUser, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctxUser)
	var usersFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &usersFound); err != nil {
		return nil, err
	}
	return usersFound, nil
}

//retru na slice with all the users
func GetAllUsers() ([]User, error) {
	cur, err := collectionUser.Find(ctxUser, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctxUser)
	var usersFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &usersFound); err != nil {
		return nil, err
	}
	return usersFound, nil
}

//will add the user to database, return the id if succeded adding the user
func AddUser(user, pass string, authLvl int) (string, error) {
	//check if strings are empty and authlvl between 0 and 2
	if user == "" && pass == "" && (authLvl <= -1 || authLvl >= 3) {
		return "", errors.New("uncorrect/uncomplete credentials to create the user")
	}
	//check if not already registered
	found, err := IsCorrect(user, pass)

	if !found.IsStructureEmpty() {
		return "", errors.New("user already exist")
	}

	//adding user to database
	// toInsert := User{user, pass, authlvl}
	toInsert := struct {
		User  string `bson:"user, omitempty"      json: "user, omitempty"`
		Pass  string `bson:"password, omitempty"  json: "password, omitempty"`
		Level int    `bson:"authLevel, omitempty" json: "authLevel, omitempty"`
	}{
		user,
		pass,
		authLvl,
	}
	result, err := collectionUser.InsertOne(ctxUser, toInsert)
	if err != nil {
		return "", err
	}
	InsertedID := CleanMongoId(fmt.Sprintf("%v", result.InsertedID))
	return InsertedID, nil
}

func UpdateUser(user, pass string, update bson.M) error {
	_, err := collectionUser.UpdateOne(
		ctxUser,
		bson.M{"user": user, "password": pass},
		bson.D{
			{"$set", update},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(user, pass string) error {
	_, err := collectionUser.DeleteOne(ctxUser, bson.M{"user": user, "password": pass})
	if err != nil {
		return err
	}
	return nil
}

func testAggregate() ([]User, error) {
	names := []string{"vano", "MoraGames"}
	// query := bson.D{{"$match", bson.D{{"user", bson.M{"$in": names}}}}}
	query := bson.D{{"$match", bson.D{{"user", bson.M{"$in": names}}}}}
	sort := bson.D{{"$sort", bson.M{"authLevel": 1}}}
	cur, err := collectionUser.Aggregate(ctxUser, mongo.Pipeline{query, sort})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctxUser)
	var usersFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &usersFound); err != nil {
		return nil, err
	}
	return usersFound, nil
}

func init() {
	ConnectToDatabaseUsers()
}

// func main() {
// 	ConnectToDatabaseUsers()
// 	fmt.Println(testAggregate())
// }

/*run this main to see all functionality
func main() {
	// Check the connection
	err := ConnectToDatabaseUsers()
	if err != nil {
		fmt.Println(err)
	}


	// users, _ := GetAllUsers()
	// fmt.Println("all users", users)
	// fmt.Println()

	// query := bson.M{"user": "cami<3"}
	// users, _ = QueryUsers(query)
	// fmt.Println("user == cami", users)
	// fmt.Println()

	// update := bson.M{"user": "cami<3"}
	// UpdateUser("cami", "HelloThere:D123!!!", update)
	// users, _ = GetAllUsers()
	// fmt.Println("updated cami into cami<3", users)
	// fmt.Println()

	// DeleteUser("ciao", "camiCwute")
	// users, _ = GetAllUsers()
	// fmt.Println("deleted ciao", users)
	// fmt.Println()

	// found, err := IsCorrect("vano", "HelloThere:D123!!!")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("correct:", found)

	// found, err = IsCorrect("cami", "HelloThere:D123")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("incorrect psw:", found)

	// found, err = IsCorrect("chonky", "HelloThere:D123!!!")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("incorrect user:", found)
}
*/
