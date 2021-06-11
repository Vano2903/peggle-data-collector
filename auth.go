package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	ctx        context.Context
	collection *mongo.Collection
)

type User struct {
	ID    primitive.ObjectID `bson:"_id, omitempty"       json: "_id, omitempty"`
	User  string             `bson:"user, omitempty"      json: "user, omitempty"`
	Pass  string             `bson:"password, omitempty"  json: "password, omitempty"`
	Level int                `bson:"authLevel, omitempty" json: "authLevel, omitempty"` //0 admin (every priviliges), 1 every power relative to the games, 2 just adding
}

//will connect to database on user's collection
func connectToDatabaseUsers() error {
	//get context
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.9:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	//check if connection is established
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	//assign to the global variable "collection" the users' collection
	collection = client.Database("qdss-peggle").Collection("users")
	return nil
}

//check if user exist in database and will return -1 if not found on the other hand will return the authorization level
//? I belive there is a bettere way to do this but rn i dont really know
func IsCorrect(user, pass string) (int, error) {
	//search in database
	cur, err := collection.Find(ctx, bson.M{"user": user, "password": pass})
	if err != nil {
		return -1, err
	}
	defer cur.Close(ctx)
	var userFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &userFound); err != nil {
		return -1, err
	}

	//check if user exist
	if len(userFound) != 0 {
		if userFound[0].User == user && userFound[0].Pass == pass {
			return userFound[0].Level, nil
		}
		return -1, nil
	} else {
		return -1, nil
	}
}

//return a []User given a bson.M query (possible queries: user, authLevel)
func QueryUsers(query bson.M) ([]User, error) {
	if _, ok := query["password"]; ok {
		return nil, errors.New("you can't query over passwords")
	}
	cur, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var usersFound []User

	//convert cur in []User
	if err = cur.All(context.TODO(), &usersFound); err != nil {
		return nil, err
	}
	return usersFound, nil
}

//retru na slice with all the users
func GetAllUsers() ([]User, error) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
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
	lvl, err := IsCorrect(user, pass)
	if err != nil {
		return "", err
	}
	if lvl != -1 {
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
	result, err := collection.InsertOne(ctx, toInsert)
	if err != nil {
		return "", err
	}
	InsertedID := fmt.Sprintf("%v", result.InsertedID)
	InsertedID = strings.Replace(InsertedID, "ObjectID(\"", "", -1)
	InsertedID = strings.Replace(InsertedID, "\")", "", -1)
	return InsertedID, nil
}

func UpdateUser(user, pass string, update bson.M) error {
	_, err := collection.UpdateOne(
		ctx,
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
	_, err := collection.DeleteOne(ctx, bson.M{"user": user, "password": pass})
	if err != nil {
		return err
	}
	return nil
}

/*
* run this main to see all functionality
func main() {
	// Check the connection
	err := connectToDatabaseUsers()
	if err != nil {
		fmt.Println(err)
	}

	users, _ := GetAllUsers()
	fmt.Println("all users", users)
	fmt.Println()

	query := bson.M{"user": "cami"}
	users, _ = QueryUsers(query)
	fmt.Println("user == cami", users)
	fmt.Println()

	update := bson.M{"user": "cami<3"}
	UpdateUser("cami", "HelloThere:D123!!!", update)
	users, _ = GetAllUsers()
	fmt.Println("updated cami into cami<3", users)

	AddUser("ciao", "camiCwute", 1)
	users, _ = GetAllUsers()
	fmt.Println("added user: ciao, camiCwute, 1", users)
	fmt.Println()

	DeleteUser("ciao", "camiCwute")
	users, _ = GetAllUsers()
	fmt.Println("deleted ciao", users)
	fmt.Println()

	found, err := IsCorrect("vano", "HelloThere:D123!!!")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("correct:", found)

	found, err = IsCorrect("cami", "HelloThere:D123")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("incorrect psw:", found)

	found, err = IsCorrect("chonky", "HelloThere:D123!!!")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("incorrect user:", found)
}
*/
