package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientStats     *mongo.Client
	ctxStats        context.Context
	collectionStats *mongo.Collection
)

type OverallStats struct {
	Id      int          `bson:"id" json:"-"`                                  //used just to update the stats in the database
	Generic GenericStats `bson:"generic, omitempty" json:"generic, omitempty"` //generic's statistics
	Synergo PlayerStats  `bson:"synergo, omitempty" json:"synergo, omitempty"` //synergo's statistics
	Redez   PlayerStats  `bson:"redez, omitempty" json:"redez, omitempty"`     //redez's statistics
}

type GenericStats struct {
	TotalTimeWatched    int      `bson:"totTimeWatched, omitempty" json:"totTimeWatched, omitempty"`       //total ammount of seconds watched
	TotalEpisodesStored int      `bson:"totEpisodesStored, omitempty" json:"totEpisodesStored, omitempty"` //number of all the episodes stored
	Collaborators       []string `bson:"collaborators, omitempty" json:"collaborators, omitempty"`         //number of all the users that have at least 1 commit
}

//"method of genericStats" check if the collaborators is already stored in []string collaborators
func (g GenericStats) containsUser(u string) bool {
	for _, col := range g.Collaborators {
		if col == u {
			return true
		}
	}
	return false
}

type PlayerStats struct {
	TotalPoints int       `bson:"totPoints, omitempty" json:"totPoints, omitempty"` //total ammount of points made by a player
	TotalN25    int       `bson:"totn25, omitempty" json:"totn25, omitempty"`       //total ammount of -25% made by a player
	TotalWins   int       `bson:"totalWins, omitempty" json:"totWins, omitempty"`   //total ammount of wins made by a player
	FEstats     FEstats   `bson:"FEstats, omitempty" json:"FEstats, omitempty"`     //extreme fever stats of a player
	ChartStats  CharStats `bson:"charStats, omitempty" json:"charStats,omitempty"`  //character stats of a player
}

type FEstats struct {
	N5000         int `bson:"n5000" json:"n5000"`                 //total ammount of 5000 type extreme fever
	N25000        int `bson:"n25000" json:"n25000"`               //total ammount of 25000 type extreme fever
	N50000        int `bson:"n50000" json:"n50000"`               //total ammount of 50000 type extreme fever
	TotPointsMade int `bson:"totPointsMade" json:"totPointsMade"` //sum of all the points made with extreme fevers
}

type CharStats struct {
	Cas int `bson:"cas" json:"cas"` //total ammount of times the player choose this character
	Uni int `bson:"uni" json:"uni"` //total ammount of times the player choose this character
	Zuc int `bson:"zuc" json:"zuc"` //total ammount of times the player choose this character
	Gat int `bson:"gat" json:"gat"` //total ammount of times the player choose this character
	Ali int `bson:"ali" json:"ali"` //total ammount of times the player choose this character
	Gra int `bson:"gra" json:"gra"` //total ammount of times the player choose this character
	Gir int `bson:"gir" json:"gir"` //total ammount of times the player choose this character
	Dra int `bson:"dra" json:"dra"` //total ammount of times the player choose this character
	Con int `bson:"con" json:"con"` //total ammount of times the player choose this character
	Guf int `bson:"guf" json:"guf"` //total ammount of times the player choose this character
	Sep int `bson:"sep" json:"sep"` //total ammount of times the player choose this character
}

//will connect to database on stats's collectionn
func ConnectToDatabaseStats() error {
	//get context
	ctxStats, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	//try to connect
	clientOptions := options.Client().ApplyURI(conf.Uri)
	clientStats, err := mongo.Connect(ctxStats, clientOptions)
	if err != nil {
		return err
	}

	//check if connection is established
	err = clientStats.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	//assign to the global variable "collection" the stats' collection
	collectionStats = clientStats.Database("qdss-peggle").Collection("stats")
	return nil
}

//given a game will update the extreme fever datas
func (f *FEstats) addFEData(g Game, who string) {
	if who == "s" {
		switch g.Stats.Synergo.G1.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
		switch g.Stats.Synergo.G2.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
		switch g.Stats.Synergo.G3.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
	} else {
		switch g.Stats.Redez.G1.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
		switch g.Stats.Redez.G2.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
		switch g.Stats.Redez.G3.ValFE {
		case 5000:
			f.N5000++
			f.TotPointsMade += 5000
		case 25000:
			f.N25000++
			f.TotPointsMade += 25000
		case 50000:
			f.N50000++
			f.TotPointsMade += 50000
		}
	}
}

//given a game will update the character datas
func (c *CharStats) addCharData(g Game, who string) {
	if who == "s" {
		switch g.Stats.Synergo.G1.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
		switch g.Stats.Synergo.G2.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
		switch g.Stats.Synergo.G3.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
	} else {
		switch g.Stats.Redez.G1.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
		switch g.Stats.Redez.G2.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
		switch g.Stats.Redez.G3.Character {
		case "castoro":
			c.Cas++
		case "unicorno":
			c.Uni++
		case "zucca":
			c.Zuc++
		case "gatto":
			c.Gat++
		case "alieno":
			c.Ali++
		case "granchio":
			c.Gra++
		case "girasole":
			c.Gir++
		case "drago":
			c.Dra++
		case "coniglio":
			c.Con++
		case "gufo":
			c.Guf++
		case "seppia":
			c.Sep++
		}
	}
}

//given a game will update the stats in the database
func (s OverallStats) AddStatsData(g Game) error {
	s, err := LoadStatsFromDB()
	if err != nil {
		return err
	}

	if !s.Generic.containsUser(g.AddedBy) {
		s.Generic.Collaborators = append(s.Generic.Collaborators, g.AddedBy)
	}

	s.Generic.TotalTimeWatched += g.VD.Length
	s.Generic.TotalEpisodesStored++

	s.Synergo.TotalPoints += g.Stats.Synergo.Overall.TPoints
	s.Synergo.TotalN25 += g.Stats.Synergo.Overall.T25
	s.Synergo.FEstats.addFEData(g, "s")
	s.Synergo.ChartStats.addCharData(g, "s")

	s.Redez.TotalPoints += g.Stats.Redez.Overall.TPoints
	s.Redez.TotalN25 += g.Stats.Redez.Overall.T25
	s.Redez.FEstats.addFEData(g, "r")
	s.Redez.ChartStats.addCharData(g, "r")

	if g.WonBy == 1 {
		s.Synergo.TotalWins++
	} else if g.WonBy == 0 {
		s.Redez.TotalWins++
	}
	err = s.UploadStatsToDB()
	if err != nil {
		return err
	}
	return nil
}

//this function is used for debuggin, can insert for the first time a overallStat variable in the database
func (s *OverallStats) insertFirst() error {
	_, err := collectionStats.InsertOne(ctxStats, s)
	if err != nil {
		return err
	}
	return nil
}

//update the statistics in the database
func (s *OverallStats) UploadStatsToDB() error {
	_, err := collectionStats.UpdateOne(
		ctxStats,
		bson.M{"id": 0},
		bson.D{
			{"$set", s},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//load all the statistics from the database
func LoadStatsFromDB() (OverallStats, error) {
	query := bson.M{"id": 0}
	cur, err := collectionStats.Find(ctxStats, query)
	if err != nil {
		return OverallStats{}, err
	}
	defer cur.Close(ctxStats)
	var stat []OverallStats

	//convert cur in []OverallStats
	if err = cur.All(context.TODO(), &stat); err != nil {
		return OverallStats{}, err
	}
	return stat[0], nil
}

//! use this main if you need to rewrite the whole stats db
// func main() {
// 	var queries []bson.D
// 	q := bson.D{{"$match", bson.D{{"videoData.title", bson.D{{"$regex", primitive.Regex{Pattern: "", Options: "i"}}}}}}}
// 	queries = append(queries, q)
// 	g, _ := QueryGames(queries)

// 	fmt.Println(g)

// 	var s OverallStats
// 	s.insertFirst()

// 	for _, gam := range g {
// 		s.AddStatsData(gam)
// 	}
// }
