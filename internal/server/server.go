package server

import (
	"database/sql"
	"log"
	"time"

	"github.com/abc_valera/flugo/internal/database"
	cnfg "github.com/abc_valera/flugo/internal/utils/config"
	"github.com/abc_valera/flugo/internal/utils/middleware"
	"github.com/abc_valera/flugo/internal/utils/token"
	v "github.com/abc_valera/flugo/internal/utils/validator"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Server struct {
	app *fiber.App

	config     cnfg.Config
	db         *database.Queries
	tokenMaker token.Maker
	validator  v.CustomValidator
}

func NewServer() (*Server, error) {
	s := new(Server)

	// init config
	c, err := cnfg.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	s.config = c

	// init database
	conn, err := sql.Open(s.config.DatabaseDriver, s.config.DatabaseUrl)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	s.db = database.New(conn)

	// init migrations
	m, err := migrate.New("file://internal/database/migrations", s.config.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	// ! MIGRATION DOWN FOR TEST PURPOSES !
	err = m.Down()
	if err != nil {
		log.Println(err)
	}
	err = m.Up()
	if err != nil {
		return nil, err
	}
	//!!!
	m.Close()

	// init fiber app with custom error handler
	s.app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"message": e.Message,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		},
	})

	// init tokenMaker
	s.tokenMaker, err = token.NewJWTMaker(s.config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	// init custom logger
	s.app.Use(logger.New(logger.Config{
		Format:     "${time} |${status}-${method}| ${path}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Europe/Helsinki",
	}))

	// init custom validator
	s.validator = &v.FlugoValidator{
		Validator: validator.New(),
	}

	return s, nil
}

func (s *Server) initRouter() {
	// for unauthorized user
	// static
	s.app.Static("/uploads", "./uploads")
	// users
	s.app.Post("/users", s.createUser)
	s.app.Post("/users/login", s.loginUser)
	s.app.Get("/users/verify/email", s.verifyEmail)
	s.app.Get("/users", s.listUsers)
	// jokes
	s.app.Get("/jokes", s.listJokes)
	s.app.Get("/jokes/:id", s.getJoke)
	s.app.Get("/jokes_by/:username", s.listJokesByAuthor)

	// for authorized users
	authMiddleware := middleware.NewAuthMiddleware(s.tokenMaker)
	auth := s.app.Group("/")
	auth.Use(authMiddleware)
	// users
	auth.Get("/users/me", s.getMe)
	auth.Put("/users/password", s.updateUserPassword)
	auth.Post("/uploads/images/avatars", s.updateUserAvatar)
	auth.Put("/users/fullname", s.updateUserFullname)
	auth.Put("/users/status", s.updateUserStatus)
	auth.Put("/users/bio", s.updateUserBio)
	auth.Delete("/users", s.deleteUser)
	// jokes
	auth.Post("/jokes", s.createJoke)
	auth.Put("/jokes/title/:id", s.updateJokeTitle)
	auth.Put("/jokes/text/:id", s.updateJokeText)
	auth.Put("/jokes/explanation/:id", s.updateJokeExplanation)
	auth.Delete("/jokes/:id", s.deleteJoke)
	auth.Delete("/jokes", s.deleteJokesByAuthor)

	// !DANGEROUS FUNCTION FOR TEST ONLY!
	s.app.Delete("/users_ALL", s.deleteAllUsers)
	s.app.Delete("/jokes_ALL", s.deleteAllJokes)
}

func (s *Server) Start() {
	s.initRouter()
	s.app.Listen(s.config.PORT)
}
