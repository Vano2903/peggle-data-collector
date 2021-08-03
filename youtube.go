package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

type VideoData struct {
	Id             string             `bson:"id, omitempty" json:"id,omitempty"`
	Title          string             `bson:"title, omitempty" json:"title, omitempty"`
	ThumbMaxResUrl string             `bson:"thumbMaxResUrl, omitempty" json:"thumbMaxResUrl, omitempty"`
	UploadDate     primitive.DateTime `bson:"uploadDate, omitempty" json:"uploadDate, omitempty"`
	Length         int                `bson:"length, omitempty`
}

//TODO
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
	call := service.Videos.List([]string{"snippet", "contentDetails"})
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

	length := response.Items[0].ContentDetails.Duration
	length = strings.Replace(length, "PT", "", -1)
	length = strings.Replace(length, "S", "", -1)
	lengthItems := strings.Split(length, "M")
	var lengthItemsInt [2]int
	lengthItemsInt[0], err = strconv.Atoi(lengthItems[0])
	if err != nil {
		return err
	}
	lengthItemsInt[1], err = strconv.Atoi(lengthItems[1])
	if err != nil {
		return err
	}
	v.Length = (lengthItemsInt[0] * 60) + lengthItemsInt[1]

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

// func main() {
// 	var v VideoData
// 	v.GetYoutubeDataFromId("z8bj_wLQf5I")
// 	fmt.Println(v)
// }
