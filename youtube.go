package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	Id             string `bson:"id, omitempty" json:"id,omitempty"`
	Title          string `bson:"title, omitempty" json: "title, omitempty"`
	ThumbMaxResUrl string `bson:"thumbMaxResUrl, omitempty" json: "thumbMaxResUrl, omitempty"`
	UploadDate     DateYt `bson:"uploadDate, omitempty" json: "uploadDate, omitempty"`
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
	v.Id = id
	v.Title = response.Items[0].Snippet.Title
	v.ThumbMaxResUrl = response.Items[0].Snippet.Thumbnails.Maxres.Url
	v.UploadDate.ParseString(response.Items[0].Snippet.PublishedAt)
	return nil
}

type DateYt struct {
	Day   int `bson:"day, omitempty" json: "day, omitempty"`
	Month int `bson:"month, omitempty" json: "month, omitempty"`
	Year  int `bson:"year, omitempty" json: "year, omitempty"`
}

//convert the way youtube store dates in DateYt struct
func (v *DateYt) ParseString(ytDate string) error {
	da := ytDate[:len(ytDate)-10]
	ele := strings.Split(da, "-")

	y, err := strconv.Atoi(ele[0])
	if err != nil {
		return err

	}
	m, err := strconv.Atoi(ele[1])
	if err != nil {
		return err
	}
	d, err := strconv.Atoi(ele[2])
	if err != nil {
		return err
	}

	v.Year = y
	v.Month = m
	v.Day = d
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