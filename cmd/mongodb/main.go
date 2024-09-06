package main

import (
	"http-mongodb-api/config"
	"http-mongodb-api/pkg/db"
	"http-mongodb-api/routes"
)

func main() {

	config.InitConfig()

	db.InitMongoDb()
	defer db.CloseMongoDbClient()

	routes.Api.InitRouter()
}
