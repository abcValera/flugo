package server

import (
	"database/sql"
	"log"
	"time"

	cnfg "github.com/abc_valera/flugo/internal/config"
	"github.com/abc_valera/flugo/internal/database"
	"github.com/abc_valera/flugo/internal/token"
	fv "github.com/abc_valera/flugo/internal/validator"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	app *fiber.App

	db         *database.Queries
	config     cnfg.Config
	tokenMaker token.Maker
	validate   fv.CustomValidator
)

func initConfig() {
	c, err := cnfg.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	config = c
}

func initDatabase() {
	conn, err := sql.Open(config.DatabaseDriver, config.DatabaseUrl)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal("connection is aborted: ", err)
	}
	db = database.New(conn)
}

func initMigration() {
	m, err := migrate.New("file://internal/database/migrations", config.DatabaseUrl)
	if err != nil {
		log.Fatal("migration initialize failed: ", err)
	}
	// ! MIGRATION DOWN FOR TEST PURPOSES !
	err = m.Down()
	if err != nil {
		log.Println("migration down failed: ", err)
	}
	err = m.Up()
	if err != nil {
		log.Println("migration up failed: ", err)
	}
}

func initServer() {
	// returns instance of fiber.App with a custom error handler
	app = fiber.New(fiber.Config{
		ErrorHandler: initCustomErrorHandler(),
	})

	initTokenMaker()
	initRequestLogger()
	initValidator()
	initRoutes()
}

func initCustomErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		if e, ok := err.(*fiber.Error); ok {
			return c.Status(e.Code).JSON(fiber.Map{
				"message": e.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
}

func initTokenMaker() {
	tm, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: " + err.Error())
	}
	tokenMaker = tm
}

func initRequestLogger() {
	app.Use(logger.New(logger.Config{
		Format:     "${time} |${status}-${method}| ${path}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Europe/Helsinki",
	}))
}

func initValidator() {
	validate = &fv.FlugoValidator{Validator: validator.New()}
}
