package main

type endpoints string

const (

	//base route
	root endpoints = "/"

	//API route
	endApi endpoints = "/api/"

	//users endpoint
	usersLogin endpoints = "/users/login"
)

func (e endpoints) String() string {
	return string(e)
}
