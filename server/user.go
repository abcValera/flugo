package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/abc_valera/flugo/database"
	"github.com/abc_valera/flugo/middleware"
	"github.com/abc_valera/flugo/token"
	"github.com/abc_valera/flugo/utils"
	"github.com/gofiber/fiber/v2"
)

// UserResponse type is returned back with response. It omits unnecessary data from the database's user type.
type userResponse struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
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

func createUser(c *fiber.Ctx) error {
	req := new(createUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := db.CreateUser(c.Context(), database.CreateUserParams{
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

func loginUser(c *fiber.Ctx) error {
	req := new(loginUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	accessToken, err := tokenMaker.CreateToken(user.ID, user.Username, user.Email, config.AccessTokenDuration)
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

func verifyEmail(c *fiber.Ctx) error {
	req := new(verifyEmailRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusOK, "Email is not registered yet")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return fiber.NewError(fiber.StatusBadRequest, "Email is registered already")
}

func listUsers(c *fiber.Ctx) error {
	first, err := strconv.Atoi(c.Query("first"))
	if err != nil {
		fiber.NewError(http.StatusBadRequest, err.Error())
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		fiber.NewError(http.StatusBadRequest, err.Error())
	}

	users, err := db.ListUsers(c.Context(), database.ListUsersParams{
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

func getMe(c *fiber.Ctx) error {
	user, err := db.GetUserByID(c.Context(), c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID)
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

func updateUserPassword(c *fiber.Ctx) error {
	req := new(updateUserPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sessID := c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID
	oldUser, err := db.GetUserByID(c.Context(), sessID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := utils.CheckPassword(req.OldPassword, oldUser.HashedPassword); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	user, err := db.UpdateUserPassword(c.Context(), database.UpdateUserPasswordParams{
		ID:             int32(sessID),
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusBadRequest).JSON(newUserResponse(user))
}

type updateUserFullnameRequest struct {
	Fullname string `json:"fullname" validate:"required"`
}

func updateUserFullname(c *fiber.Ctx) error {
	req := new(updateUserFullnameRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := db.UpdateUserFullname(c.Context(), database.UpdateUserFullnameParams{
		ID:       c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID,
		Fullname: req.Fullname,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusBadRequest).JSON(newUserResponse(user))
}

type updateUserStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func updateUserStatus(c *fiber.Ctx) error {
	req := new(updateUserStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := db.UpdateUserStatus(c.Context(), database.UpdateUserStatusParams{
		ID:     c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID,
		Status: req.Status,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUserResponse(user))
}

type updateUserBioRequest struct {
	Bio string `json:"bio" validate:"required"`
}

func updateUserBio(c *fiber.Ctx) error {
	req := new(updateUserBioRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := db.UpdateUserBio(c.Context(), database.UpdateUserBioParams{
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

func deleteUser(c *fiber.Ctx) error {
	req := new(deleteUserRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sessID := c.Locals(middleware.AuthPayloadKey).(*token.Payload).UserID
	user, err := db.GetUserByID(c.Context(), int32(sessID))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := utils.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	err = db.DeleteJokesByAuthor(c.Context(), user.Username)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	err = db.DeleteUser(c.Context(), int32(sessID))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
