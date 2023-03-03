package server

import "github.com/abc_valera/flugo/internal/middleware"

func initRoutes() {
	// for unauthorized user
	// static
	app.Static("/uploads", "./uploads")
	// users
	app.Post("/users", createUser)
	app.Post("/users/login", loginUser)
	app.Get("/users/verify/email", verifyEmail)
	app.Get("/users", listUsers)
	// jokes
	app.Get("/jokes", listJokes)
	app.Get("/jokes/:id", getJoke)
	app.Get("/jokes_by/:username", listJokesByAuthor)

	// for authorized users
	authMiddleware := middleware.NewAuthMiddleware(tokenMaker)
	auth := app.Group("/")
	auth.Use(authMiddleware)
	// users
	auth.Get("/users/me", getMe)
	auth.Put("/users/password", updateUserPassword)
	auth.Post("/uploads/images/avatars", updateUserAvatar)
	auth.Put("/users/fullname", updateUserFullname)
	auth.Put("/users/status", updateUserStatus)
	auth.Put("/users/bio", updateUserBio)
	auth.Delete("/users", deleteUser)
	// jokes
	auth.Post("/jokes", createJoke)
	auth.Put("/jokes/title/:id", updateJokeTitle)
	auth.Put("/jokes/text/:id", updateJokeText)
	auth.Put("/jokes/explanation/:id", updateJokeExplanation)
	auth.Delete("/jokes/:id", deleteJoke)
	auth.Delete("/jokes", deleteJokesByAuthor)

	// !DANGEROUS FUNCTION FOR TEST ONLY!
	app.Delete("/users_ALL", deleteAllUsers)
	app.Delete("/jokes_ALL", deleteAllJokes)
}
