package main

type endpoints string

const (

	//*base route
	root endpoints = "/" //root enpoint, is for the homepage
	//endpoint for the single page game ({id} is the id of the game
	//unless if it's "api", "stats", "suport", "404")
	pages endpoints = "/{id}"

	//TODO API route
	//endApi endpoints = "/api/"

	//*statics
	statics endpoints = "/static/" //endpoint for the statics document (js, css, html, images)

	//*users endpoint
	usersLogin        endpoints = "/users/login"          //endpoint for the login page or to check the login informations
	usersPfp          endpoints = "/users/pfp/{user}"     //endpoint to get the profile pictures of a user ({user} is the username of the user)
	userCustomization endpoints = "/users/customization/" //endpoint for letting the user customise his account (like the profile picture or the password)

	//*commits endpoints

	//endpoint to handle the commits ({param} can be:
	// totCommits to get the number of commits,
	//years to get all the years a user has at least a commit in,
	//year get all the commit in the required year,
	//add is for adding a commit in the user stats,
	//anything else will return bad request)
	getCommits endpoints = "/commit/{param}"

	//*games endpoints
	games      endpoints = "/games/search"      //endpoint for quering a game using url query
	checkGame  endpoints = "/games/check/{id}"  //endpoint to check if an id is on the database ({id} is the id of the game)
	addGame    endpoints = "/games/add"         //endpoint which the user can send the game json
	updateGame endpoints = "/games/update/{id}" //endpoint used to update a game with a new game json
	deleteGame endpoints = "/games/delete/{id}" //endpoint to delete from database a game with id

	//*statistic endpoints
	stats endpoints = "/stats/{id}" //endpoint to get the stats data, ({id} can be all, generic, synergo, redez)
)

//convert endpoint to string
func (e endpoints) String() string {
	return string(e)
}
