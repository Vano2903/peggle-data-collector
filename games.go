package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientGame     *mongo.Client
	ctxGame        context.Context
	collectionGame *mongo.Collection
)

type Game struct {
	// ID      primitive.ObjectID `bson:"_id" json:"-"`
	VD      VideoData `bson:"videoData" json:"videoData"` //videoData rappresent the information of the video
	WonBy   int       `bson:"wonBy" json:"wonBy"`         //syn = 1, red = 0, pareggio/null/whatever = -1
	Stats   Players   `bson:"stats" json:"stats"`         //statistics of synergo and redez
	Comment string    `bson:"comment" json:"comment"`     //comment from the user
	AddedBy string    `bson:"addedBy" json:"addedBy"`     //usename of who added the game
}

type Players struct {
	Synergo Player `bson:"synergo" json:"synergo"`
	Redez   Player `bson:"redez" json:"redez"`
}

type Player struct {
	Overall Overall   `bson:"overall" json:"overall"` //overall data
	G1      GameStats `bson:"g1" json:"g1"`           //game 1 data
	G2      GameStats `bson:"g2" json:"g2"`           //game 2 data
	G3      GameStats `bson:"g3" json:"g3"`           //game 3 data
}

type Overall struct {
	TPoints int `bson:"tPoints" json:"tPoints"` //number of points made in a game by a player
	T25     int `bson:"t25" json:"t25"`         //number of -25% made in a game by a player
}

type GameStats struct {
	Points    int    `bson:"points" json:"points"`       //points made in a single part of the game by a player
	N25       int    `bson:"n25" json:"n25"`             //-25% made in a single part of the game by a player
	ValFE     int    `bson:"valFe" json:"valFe"`         //extreme fever made in a single part of the game by a player
	Character string `bson:"character" json:"character"` //character choosen in a single part of the game by a player
}

//will connect to database on games's collectionn
func ConnectToDatabaseGame() error {
	//get context
	ctxGame, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI(conf.Uri)
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

//will query a game using aggregate type for querying
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

//check if the id exist in the database, true if found
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

//TODO check if toAdd is not completed
//given a game this function will add the game in the database
func AddGame(toAdd Game) (string, error) {
	//check if not already stored
	found, _ := CheckIfExist(toAdd.VD.Id)
	if found == true {
		return "", errors.New("game already stored, if you want to update an already stored game use UpdateGame")
	}

	//adding game to database
	toAddNoId := struct {
		VD      VideoData `bson:"videoData" json:"videoData"`
		WonBy   int       `bson:"wonBy" json:"wonBy"` //syn = 1, red = 0, pareggio/null/whatever = -1
		Stats   Players   `bson:"stats" json:"stats"`
		Comment string    `bson:"comment" json:"comment"`
		AddedBy string    `bson:"addedBy" json:"addedBy"`
	}{
		toAdd.VD,
		toAdd.WonBy,
		toAdd.Stats,
		toAdd.Comment,
		toAdd.AddedBy,
	}
	result, err := collectionGame.InsertOne(ctxGame, toAddNoId)
	if err != nil {
		return "", err
	}
	InsertedID := CleanMongoId(fmt.Sprintf("%v", result.InsertedID))
	return InsertedID, nil
}

//update just a section of a game
func PartialUpdateGame(id string, update bson.M) error {
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

//update completly a game
func FullUpdateGame(id string, update Game) error {
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

//delete a game
func DeleteGame(id string) error {
	_, err := collectionGame.DeleteOne(ctxGame, bson.M{"videoData.id": id})
	if err != nil {
		return err
	}
	return nil
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
// 	game.Stats.Synergo.G1 = GameStats{34195, 1, 0, "puttana eva"}
// 	game.Stats.Synergo.G2 = GameStats{39830, 1, 0, "girasole"}
// 	game.Stats.Synergo.G3 = GameStats{58390, 1, 0, "gatto"}

// 	game.Stats.Redez.Overall = Overall{114505, 2}
// 	game.Stats.Redez.G1 = GameStats{32860, 0, 0, "castoro"}
// 	game.Stats.Redez.G2 = GameStats{42840, 1, 5000, "alieno"}
// 	game.Stats.Redez.G3 = GameStats{38805, 1, 0, "zucca"}
// 	// ConnectToDatabaseUsers()
// 	fmt.Println(FullUpdateGame("IwvS8ft7DM8", game))

// fmt.Println(AddGame(game))
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
