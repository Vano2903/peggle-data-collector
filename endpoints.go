package main

type endpoints string

const (

	//base route
	root endpoints = "/"

	//API route
	endApi endpoints = "/api/"

	//statics
	statics endpoints = "/static/"

	//users endpoint
	usersLogin     endpoints = "/users/login"
	getCommits     endpoints = "/commit/{param}"
	getCommitsYear endpoints = "/commit/{year}"
)

func (e endpoints) String() string {
	return string(e)
}
