package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"gopkg.in/yaml.v2"
)

var (
	conf config
)

type config struct {
	Apikey string `yaml:"yt_api_v3_key"`
}

type VideoData struct {
	Id             string             `bson:"id, omitempty" json:"id,omitempty"`
	Title          string             `bson:"title, omitempty" json:"title, omitempty"`
	ThumbMaxResUrl string             `bson:"thumbMaxResUrl, omitempty" json:"thumbMaxResUrl, omitempty"`
	UploadDate     primitive.DateTime `bson:"uploadDate, omitempty" json:"uploadDate, omitempty"`
}

func (v *VideoData) CheckIfNotCompleted() bool {
	if v.Title == "" || v.ThumbMaxResUrl == "" {
		return false
	}
	return true
}

//fill the structure with all the data given the id of the video
func (v *VideoData) GetYoutubeDataFromId(id string) error {
	service, err := GetYoutubeService(conf.Apikey)
	if err != nil {
		return err
	}
	call := service.Videos.List([]string{"snippet"})
	call.Id(id)

	response, err := call.Do()
	if err != nil {
		return err
	}
	if len(response.Items) <= 0 {
		return errors.New("no video found, check the id")
	}
	fmt.Println(response.Items[0].Snippet)
	v.Id = id
	v.Title = response.Items[0].Snippet.Title

	if response.Items[0].Snippet.Thumbnails.Maxres != nil {
		v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.Maxres.Url
	} else if response.Items[0].Snippet.Thumbnails.Standard != nil {
		v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.Standard.Url
	} else if response.Items[0].Snippet.Thumbnails.High != nil {
		v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.High.Url
	} else if response.Items[0].Snippet.Thumbnails.Medium != nil {
		v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.Medium.Url
	} else if response.Items[0].Snippet.Thumbnails.Default.Url != "" {
		v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.Default.Url
	} else {
		v.ThumbMaxResUrl = ""
	}

	upDate, err := time.Parse(time.RFC3339, response.Items[0].Snippet.PublishedAt)
	if err != nil {
		return err
	}
	v.UploadDate = primitive.NewDateTimeFromTime(upDate)
	return nil
}

//return the youtube service given a valid youtube api key
func GetYoutubeService(key string) (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: key},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}

//get the api key from config.yaml
func init() {
	dat, err := ioutil.ReadFile("config.yaml")
	err = yaml.Unmarshal([]byte(dat), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// func main() {
// 	var v VideoData
// 	v.GetYoutubeDataFromId("S0-4ouN35gw")
// 	fmt.Println(v)
// }
