package main

type endpoints string

const (

	//base route
	root endpoints = "/"

	//API route
	endApi endpoints = "/api/"

	//users endpoint
	endAdmin            endpoints = "/user/"
	adminAddLearning    endpoints = "/learning/add"
	adminGetLearning    endpoints = "/learning/get"
	adminSetLearning    endpoints = "/learning/set"
	adminModifyLearning endpoints = "/learning/modify"
)

func (e endpoints) String() string {
	return string(e)
}
