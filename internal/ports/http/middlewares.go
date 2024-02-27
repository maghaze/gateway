package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *Server) optionalAuthentication(c *fiber.Ctx) error {
	headerBytes := c.Request().Header.Peek("Authorization")
	header := strings.TrimPrefix(string(headerBytes), "Bearer ")

	if len(header) == 0 {
		return c.Next()
	}

	id, err := s.auth.Authenticate(c.Context(), header)
	if err != nil {
		s.logger.Error("Invalid token header", zap.Error(err))
		return c.Next()
	}

	c.Request().Header.Add("X-User-Id", strconv.FormatUint(id, 10))
	c.Request().Header.Del("Authorization")

	return c.Next()
}

// requiredAuthentication will extract the token and put the user-information
// to the request header (before it being redirected)
func (s *Server) requiredAuthentication(c *fiber.Ctx) error {
	header := c.Request().Header.Peek("Authorization")

	if len(header) == 0 {
		s.logger.Error("Missing authorization header")
		response := "please provide your authentication information"
		return c.Status(http.StatusUnauthorized).SendString(response)
	}

	id, err := s.auth.Authenticate(c.Context(), string(header))
	if err != nil {
		s.logger.Error("Invalid token header", zap.Error(err))
		response := "invalid token header, please login again"
		return c.Status(http.StatusUnauthorized).SendString(response)
	}

	c.Request().Header.Add("X-User-Id", strconv.FormatUint(id, 10))
	c.Request().Header.Del("Authorization")

	return c.Next()
}
