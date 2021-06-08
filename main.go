package main

import (
	"log"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	home, err := os.ReadFile("collector-page/data.html")
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(home)
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", nil)
}
