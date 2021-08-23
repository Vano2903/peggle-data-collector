package main

type endpoints string

const (

	//base route
	root  endpoints = "/"
	pages endpoints = "/{id}"

	//API route
	// endApi endpoints = "/api/"

	//statics
	statics endpoints = "/static/"

	//users endpoint
	usersLogin        endpoints = "/users/login"
	usersPfp          endpoints = "/users/pfp/{user}"
	userCustomization endpoints = "/users/customization/"

	//commits endpoints
	getCommits endpoints = "/commit/{param}"

	//games endpoints
	games      endpoints = "/games/search"
	checkGame  endpoints = "/games/check/{id}"
	addGame    endpoints = "/games/add"
	updateGame endpoints = "/games/update/{id}"
	deleteGame endpoints = "/games/delete/{id}"

	//statistic endpoints
	stats endpoints = "/stats/{id}"
)

func (e endpoints) String() string {
	return string(e)
}
