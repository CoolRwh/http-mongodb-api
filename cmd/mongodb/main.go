package main

import (
	"http-mongodb-api/db"
	"http-mongodb-api/routes"
)

func main() {

	db.InitMongoDb()

	routes.Api.InitRouter()

}
