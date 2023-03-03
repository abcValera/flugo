package server

func Start() {
	initConfig()

	initDatabase()
	initMigration()

	initServer()

	app.Listen(config.PORT)
}
