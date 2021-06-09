package main

type endpoints string

const (

	//Route di base
	root endpoints = "/"

	//API route
	endApi endpoints = "/api/"

	//API admin
	endAdmin            endpoints = "/admin/"
	adminAddLearning    endpoints = "/learning/add"
	adminGetLearning    endpoints = "/learning/get"
	adminSetLearning    endpoints = "/learning/set"
	adminModifyLearning endpoints = "/learning/modify"
)

func (e endpoints) String() string {
	return string(e)
}
