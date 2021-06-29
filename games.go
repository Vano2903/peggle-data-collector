package main

import (
	"context"
	"time"

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
	ID        primitive.ObjectID `bson:"_id, omitempty" json:"id,omitempty"`
	VideoData VideoData          `bson:"videoData, omitempty" json:"videoData, omitempty"`
	WonBy     string             `bson:"wonBy, omitempty" json:"wonBy,omitempty"`
	Stats     Players            `bson:"stats, omitempty" json:"stats,omitempty"`
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

//will connect to database on games's collectionn
func ConnectToDatabaseGame() error {
	//get context
	ctxGame, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.9:27017")
	clientGame, err := mongo.Connect(ctxUser, clientOptions)
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
