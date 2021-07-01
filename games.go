package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientGame     *mongo.Client
	ctxGame        context.Context
	collectionGame *mongo.Collection
)

type Game struct {
	ID      primitive.ObjectID `bson:"_id, omitempty" json:"id,omitempty"`
	VD      VideoData          `bson:"videoData, omitempty" json:"videoData, omitempty"`
	WonBy   int                `bson:"wonBy, omitempty" json:"wonBy,omitempty"` //syn = 1, red = 0, pareggio/null/whatever = -1
	Stats   Players            `bson:"stats, omitempty" json:"stats,omitempty"`
	Comment string             `bson:"comment, omitempty" json:"comment, omitempty"`
}

type Players struct {
	Synergo Player `bson:"synergo, omitempty" json:"synergo,omitempty"`
	Redez   Player `bson:"redez, omitempty" json:"redez,omitempty"`
}

type Player struct {
	Overall Overall   `bson:"overall, omitempty" json:"overall,omitempty"`
	G1      GameStats `bson:"g1, omitempty" json:"g1,omitempty"`
	G2      GameStats `bson:"g2, omitempty" json:"g2,omitempty"`
	G3      GameStats `bson:"g3, omitempty" json:"g3,omitempty"`
}

type Overall struct {
	TPoints int `bson:"tPoints, omitempty" json:"tPoints,omitempty"`
	T25     int `bson:"t-25, omitempty" json:"t-25,omitempty"`
}

type GameStats struct {
	Points    int    `bson:"points, omitempty" json:"points,omitempty"`
	N25       int    `bson:"n-25, omitempty" json:"n-25,omitempty"`
	ValFE     int    `bson:"valFE, omitempty" json:"valFE,omitempty"`
	Character string `bson:"character, omitempty" json:"character,omitempty"`
}

//*
//will connect to database on games's collectionn
func ConnectToDatabaseGame() error {
	//get context
	ctxGame, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.9:27017")
	clientGame, err := mongo.Connect(ctxGame, clientOptions)
	if err != nil {
		return err
	}

	//check if connection is established
	err = clientGame.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	//assign to the global variable "collection" the games' collection
	collectionGame = clientGame.Database("qdss-peggle").Collection("games")
	return nil
}

//*
func QueryGames(q []bson.D) ([]Game, error) {
	cur, err := collectionGame.Aggregate(ctxGame, q)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctxGame)
	var gamesFound []Game

	//convert cur in []Game
	if err = cur.All(context.TODO(), &gamesFound); err != nil {
		return nil, err
	}
	return gamesFound, nil
}

//*
func CheckIfExist(id string) (bool, error) {
	//search in database
	cur, err := collectionGame.Find(ctxGame, bson.M{"videoData.id": id})
	if err != nil {
		return true, err
	}
	defer cur.Close(ctxGame)
	var gamesFound []Game

	//convert cur in []User
	if err = cur.All(context.TODO(), &gamesFound); err != nil {
		return true, err
	}
	if len(gamesFound) > 0 {
		return true, errors.New("game already stored")
	}
	return false, nil
}

//*
func AddGame(toAdd Game) (string, error) {
	//TODO check if toAdd is not completed

	//check if not already stored
	found, err := CheckIfExist(toAdd.VD.Id)
	if err != nil {
		return "", errors.New(err.Error() + "if you want to update an already stored game use UpdateGame")
	}
	if found == true {
		return "", errors.New("game already stored, ")
	}

	//adding game to database
	toAddNoId := struct {
		VD    VideoData `bson:"videoData, omitempty" json:"videoData, omitempty"`
		WonBy int       `bson:"wonBy, omitempty" json:"wonBy,omitempty"` //syn = 1, red = 0, pareggio/null/whatever = -1
		Stats Players   `bson:"stats, omitempty" json:"stats,omitempty"`
	}{
		toAdd.VD,
		toAdd.WonBy,
		toAdd.Stats,
	}
	fmt.Println(toAddNoId)
	result, err := collectionGame.InsertOne(ctxGame, toAddNoId)
	if err != nil {
		return "", err
	}
	InsertedID := CleanMongoId(fmt.Sprintf("%v", result.InsertedID))
	return InsertedID, nil
}

//*
func UpdateGame(id string, update bson.M) error {
	_, err := collectionGame.UpdateOne(
		ctxGame,
		bson.M{"videoData.id": id},
		bson.D{
			{"$set", update},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//*
func DeleteGame(id string) error {
	_, err := collectionGame.DeleteOne(ctxGame, bson.M{"videoData.id": id})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	ConnectToDatabaseGame()
}

// func main() {
// 	var game Game
// 	err := game.VD.GetYoutubeDataFromId("IwvS8ft7DM8")
// 	if err != nil {
// 		panic(err)
// 	}
// 	game.WonBy = 0
// 	//syn
// 	game.Stats.Synergo.Overall.TPoints = 132415
// 	game.Stats.Synergo.Overall.T25 = 3
// 	game.Stats.Synergo.G1 = GameStats{34195, 1, 0, "unicorno"}
// 	game.Stats.Synergo.G2 = GameStats{39830, 1, 0, "girasole"}
// 	game.Stats.Synergo.G3 = GameStats{58390, 1, 0, "gatto"}

// 	game.Stats.Redez.Overall = Overall{114505, 2}
// 	game.Stats.Redez.G1 = GameStats{32860, 0, 0, "castoro"}
// 	game.Stats.Redez.G2 = GameStats{42840, 1, 5000, "alieno"}
// 	game.Stats.Redez.G3 = GameStats{38805, 1, 0, "zucca"}
// 	// ConnectToDatabaseUsers()

// 	fmt.Println(AddGame(game))
// fmt.Println(QueryGame())

// update := bson.M{"wonBy": 1}
// fmt.Println(UpdateGame("IwvS8ft7DM8", update))
// fmt.Println(DeleteGame("IwvS8ft7DM9"))

// q1 := bson.D{{"$match", bson.D{{"wonBy", bson.M{"$in": []int{1}}}}}}
// q3 := bson.D{{"$match", bson.D{{"videoData.title", bson.M{"$in": []string{"PEGGLE: NON E' POSSIBILE CHE VADA  COSI"}}}}}}
// // q3 := bson.D{{"$match", bson.D{{"authLevel", bson.M{"$in": []int{0}}}}}}
// q2 := bson.D{{"$sort", bson.M{"wonBy": 1}}}

// query := []bson.D{q3, q2}

// result, err := QueryGames(query)
// fmt.Println(result)
// if err != nil {
// 	panic(err)
// }
// res, err := json.MarshalIndent(&result, "", "\t")
// if err != nil {
// 	panic(err)
// }
// fmt.Println(string(res))
// }
