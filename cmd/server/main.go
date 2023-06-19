package main

import (
	config "passKeeper/config/server"
	app "passKeeper/internal/models/app"
	db "passKeeper/internal/models/database"
)

func main() {
	sc := config.NewServerConfig()
	conn := db.ConnectDB(sc.Database)
	accountRepo := db.GetAccountRepo(conn)
	secretRepo := db.GetSecretRepo(conn)
	migrationRepo := db.GetMigrationRepo(conn)
	app := app.NewApp(*sc, accountRepo, secretRepo, migrationRepo)
	app.CreateTables()
	app.StartWebServer()

}
