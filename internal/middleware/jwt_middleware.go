package middleware

import (
	"net/http"
	"strings"

	"github.com/abc_valera/flugo/internal/token"
	"github.com/gofiber/fiber/v2"
)

const (
	AuthHeaderKey  = "authorization"
	AuthTypeBearer = "bearer"
	AuthPayloadKey = "auth_payload"
)

func NewAuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthHeaderKey)
		if len(authHeader) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "authorization is not provided")
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return fiber.NewError(http.StatusUnauthorized, "invalid authorization")
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthTypeBearer {
			return fiber.NewError(http.StatusUnauthorized, "this authorization type is not supported")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}
		c.Locals(AuthPayloadKey, payload)

		return c.Next()
	}
}
