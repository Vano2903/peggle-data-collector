package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Apikey string `yaml:"yt_api_v3_key"` //youtube api key
	Uri    string `yaml:"mongo_uri"`     //atlas mongodb uri
}

var (
	conf config
)

//printInternalErr set the status code to 500 of the http response
func PrintInternalErr(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "Internal Server Error", "error": "%s"}`, err)))
}

//printErr will return 400 error code to the client
func PrintErr(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"code": 400, "msg": "%s"}`, err)))
}

//return the unavailable service response (adding a unaviable page)
func UnavailablePage(w http.ResponseWriter) {
	//TODO add an unavailable page
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("{\"msg\": \"page unavailable at the moment\"}"))
}

//check if element is present in a slice of int
func Contains(slice []int, item int) bool {
	set := make(map[int]int, len(slice))
	for _, s := range slice {
		set[s] = 1
	}

	_, ok := set[item]
	return ok
}

//given a mongodb id as string will return just the id
func CleanMongoId(mongoId string) string {
	id := fmt.Sprintf("%v", mongoId)
	id = strings.Replace(id, "ObjectID(\"", "", -1)
	id = strings.Replace(id, "\")", "", -1)
	return id
}

//convert a slice of string (with int as string) to a slice of int
func ConvertToSliceInt(s []string) ([]int, error) {
	var converted []int
	for _, v := range s {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		converted = append(converted, i)
	}
	return converted, nil
}

//init all the connection to the database and get the api key and db uri
func init() {
	uri := os.Getenv("mongo_uri")
	yt := os.Getenv("yt_api_v3_key")
	if uri == "" || yt == "" {
		dat, err := ioutil.ReadFile("config.yaml")
		err = yaml.Unmarshal([]byte(dat), &conf)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		conf.Apikey = yt
		conf.Uri = uri
	}
	if err := ConnectToDatabaseStats(); err != nil {
		log.Fatal(err)
	}
	if err := ConnectToDatabaseGame(); err != nil {
		log.Fatal(err)
	}
	if err := ConnectToDatabaseUsers(); err != nil {
		log.Fatal(err)
	}
}
