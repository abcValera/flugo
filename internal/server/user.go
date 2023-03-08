package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abc_valera/flugo/internal/database"
	"github.com/abc_valera/flugo/internal/utils/middleware"
	"github.com/abc_valera/flugo/internal/utils/password"
	"github.com/abc_valera/flugo/internal/utils/token"

	"github.com/gofiber/fiber/v2"
)

// UserResponse type is returned back with response. It omits unnecessary data from the database's user type.
type userResponse struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Fullname  string    `json:"fullname"`
	Bio       string    `json:"bio"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Returns new UserResponse from default user type
func newUserResponse(user database.User) userResponse {
	return userResponse{
		user.ID,
		user.Username,
		user.Email,
		user.Avatar,
		user.Fullname,
		user.Bio,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	}
}

// POST REQUESTS

type createUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) createUser(c *fiber.Ctx) error {
	req := new(createUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := s.db.CreateUser(c.Context(), database.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginUserResponse struct {
	TokenType   string       `json:"token_type"`
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) loginUser(c *fiber.Ctx) error {
	req := new(loginUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := s.db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = password.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID, user.Username, user.Email, s.config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(loginUserResponse{
		TokenType:   middleware.AuthTypeBearer,
		AccessToken: accessToken,
		User:        newUserResponse(user),
	})
}

// GET REQUESTS

type verifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *Server) verifyEmail(c *fiber.Ctx) error {
	req := new(verifyEmailRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := s.db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusOK, "Email is not registered yet")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return fiber.NewError(fiber.StatusBadRequest, "Email is registered already")
}

func (s *Server) listUsers(c *fiber.Ctx) error {
	first, err := strconv.Atoi(c.Query("first"))
	if err != nil {
		fiber.NewError(http.StatusBadRequest, err.Error())
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		fiber.NewError(http.StatusBadRequest, err.Error())
	}

	users, err := s.db.ListUsers(c.Context(), database.ListUsersParams{
		Limit:  int32(size),
		Offset: int32(first),
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	usersResponse := make([]userResponse, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, newUserResponse(user))
	}

	return c.Status(fiber.StatusOK).JSON(usersResponse)
}

func (s *Server) getMe(c *fiber.Ctx) error {
	user, err := s.db.GetUserByID(c.Context(), c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(newUserResponse(user))
}

// PUT REQUESTS

type updateUserPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (s *Server) updateUserPassword(c *fiber.Ctx) error {
	req := new(updateUserPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sessID := c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID
	oldUser, err := s.db.GetUserByID(c.Context(), sessID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := password.CheckPassword(req.OldPassword, oldUser.HashedPassword); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	user, err := s.db.UpdateUserPassword(c.Context(), database.UpdateUserPasswordParams{
		ID:             int32(sessID),
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

func (s *Server) updateUserAvatar(c *fiber.Ctx) error {
	userID := c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID
	filename := fmt.Sprintf("/uploads/images/avatars/%d.png", userID)

	file, err := c.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = c.SaveFile(file, "."+filename)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := s.db.UpdateUserAvatar(c.Context(), database.UpdateUserAvatarParams{
		ID:     userID,
		Avatar: filename,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

type updateUserFullnameRequest struct {
	Fullname string `json:"fullname"`
}

func (s *Server) updateUserFullname(c *fiber.Ctx) error {
	req := new(updateUserFullnameRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := s.db.UpdateUserFullname(c.Context(), database.UpdateUserFullnameParams{
		ID:       c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID,
		Fullname: req.Fullname,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

type updateUserStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateUserStatus(c *fiber.Ctx) error {
	req := new(updateUserStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := s.db.UpdateUserStatus(c.Context(), database.UpdateUserStatusParams{
		ID:     c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID,
		Status: req.Status,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

type updateUserBioRequest struct {
	Bio string `json:"bio"`
}

func (s *Server) updateUserBio(c *fiber.Ctx) error {
	req := new(updateUserBioRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := s.db.UpdateUserBio(c.Context(), database.UpdateUserBioParams{
		ID:  c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID,
		Bio: req.Bio,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

// DELETE REQUESTS

type deleteUserRequest struct {
	Password string `json:"password" validate:"required"`
}

func (s *Server) deleteUser(c *fiber.Ctx) error {
	req := new(deleteUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := s.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sessID := c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID
	user, err := s.db.GetUserByID(c.Context(), int32(sessID))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := password.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	err = s.db.DeleteJokesByAuthor(c.Context(), user.Username)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	err = s.db.DeleteUser(c.Context(), int32(sessID))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// !DANGEROUS FUNCTION FOR TEST ONLY!
func (s *Server) deleteAllUsers(c *fiber.Ctx) error {
	err := s.db.DeleteAllJokes(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	err = s.db.DeleteAllUsers(c.Context())
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
