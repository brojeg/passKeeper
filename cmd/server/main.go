package main

import (
	config "passKeeper/config/server"
	app "passKeeper/internal/models/app"
	db "passKeeper/internal/models/database"
)

func main() {
	sc := config.NewServerConfig()
	repo := db.Get(sc.Database)
	app := app.NewApp(*sc, repo)
	app.CreateTables()
	app.StartWebServer()

}
