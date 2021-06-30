package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

//check if element is present in a slice of int
func Contains(slice []int, item int) bool {
	set := make(map[int]int, len(slice))
	for _, s := range slice {
		set[s] = 1
	}

	_, ok := set[item]
	fmt.Println(ok)
	return ok
}

func CleanMongoId(mongoId string) string {
	id := fmt.Sprintf("%v", mongoId)
	id = strings.Replace(id, "ObjectID(\"", "", -1)
	id = strings.Replace(id, "\")", "", -1)
	return id
}

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
