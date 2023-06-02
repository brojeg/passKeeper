package main

import (
	config "passKeeper/config/server"
)

func main() {
	sc := config.NewServerConfig()
	app := config.New(*sc)
	app.DefineJWTConfig()
	app.CreateTables()
	app.NewWebServer(*app)

}
