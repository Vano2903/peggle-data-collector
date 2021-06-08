package main

import "os"

var user string = "vano"
var password string = os.Getenv("PEGGLE_PASSWORD")

func isAuthorised(username, password string) bool {
	if username == user && password == password {
		return true
	}

	return false
}
