package server

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/abc_valera/flugo/internal/database"
	"github.com/abc_valera/flugo/internal/middleware"
	"github.com/abc_valera/flugo/internal/token"
	"github.com/gofiber/fiber/v2"
)

// POST REQUESTS

type createJokeRequest struct {
	Title       string `json:"title" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Explanation string `json:"explanation"`
}

func (s *Server) createJoke(c *fiber.Ctx) error {
	req := new(createJokeRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	authPayload := c.Locals(middleware.AuthPayloadKey).(*token.Payload)

	joke, err := s.db.CreateJoke(c.Context(), database.CreateJokeParams{
		Author:      authPayload.Username,
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(joke)
}

// GET REQUESTS

func (s *Server) getJoke(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	joke, err := s.db.GetJoke(c.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(joke)
}

func (s *Server) listJokesByAuthor(c *fiber.Ctx) error {
	log.Println("Here")

	queryUsername := c.Params("username")
	if queryUsername == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Provided wrong username")
	}
	first, err := strconv.Atoi(c.Query("first"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	jokes, err := s.db.ListJokesByAuthor(c.Context(), database.ListJokesByAuthorParams{
		Author: queryUsername,
		Limit:  int32(size),
		Offset: int32(first),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(jokes)
}

func (s *Server) listJokes(c *fiber.Ctx) error {
	first, err := strconv.Atoi(c.Query("first"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	jokes, err := s.db.ListJokes(c.Context(), database.ListJokesParams{
		Limit:  int32(size),
		Offset: int32(first),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(jokes)
}

// PUT REQUESTS

type updateJokeTitleRequest struct {
	Title string `json:"title"`
}

func (s *Server) updateJokeTitle(c *fiber.Ctx) error {
	req := new(updateJokeTitleRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	joke, err := s.db.UpdateJokeTitle(c.Context(), database.UpdateJokeTitleParams{
		ID:    int32(id),
		Title: req.Title,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(joke)
}

type updateJokeTextRequest struct {
	Text string `json:"text"`
}

func (s *Server) updateJokeText(c *fiber.Ctx) error {
	req := new(updateJokeTextRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	joke, err := s.db.UpdateJokeText(c.Context(), database.UpdateJokeTextParams{
		ID:   int32(id),
		Text: req.Text,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(joke)
}

type updateJokeExplanationRequest struct {
	Explanation string `json:"explanation"`
}

func (s *Server) updateJokeExplanation(c *fiber.Ctx) error {
	req := new(updateJokeExplanationRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	joke, err := s.db.UpdateJokeExplanation(c.Context(), database.UpdateJokeExplanationParams{
		ID:          int32(id),
		Explanation: req.Explanation,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(joke)
}

// DELETE REQUESTS

func (s *Server) deleteJoke(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = s.db.DeleteJoke(c.Context(), int32(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) deleteJokesByAuthor(c *fiber.Ctx) error {
	err := s.db.DeleteJokesByAuthor(c.Context(), c.Locals(middleware.AuthPayloadKey).(*token.Payload).Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// !DANGEROUS FUNCTION FOR TEST ONLY!
func (s *Server) deleteAllJokes(c *fiber.Ctx) error {
	err := s.db.DeleteAllJokes(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
