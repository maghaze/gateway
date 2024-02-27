package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"go.uber.org/zap"
)

func (handler *Server) liveness(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (handler *Server) readiness(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (s *Server) proxy(c *fiber.Ctx) error {
	path := strings.TrimPrefix(string(c.Request().URI().Path()), "/v1/")

	constructProxyURL := func(endpoint, base string) string {
		path = strings.TrimPrefix(path, endpoint)

		if len(path) > 1 {
			path = strings.TrimSuffix(path, "/")
			base += path
		}

		if query := string(c.Request().URI().QueryString()); len(query) > 0 {
			base += fmt.Sprintf("?%s", query)
		}

		return base
	}

	var location string

	if endpoint := "users"; strings.HasPrefix(path, endpoint) {
		location = constructProxyURL(endpoint, s.config.Targets.Users)
	} else if endpoint = "books"; strings.HasPrefix(path, endpoint) {
		location = constructProxyURL(endpoint, s.config.Targets.Books)
	} else {
		s.logger.Error("Invalid endpoint", zap.ByteString("URI", c.Request().URI().FullURI()))
		return c.Status(http.StatusNotFound).SendString("The requested endpoint doesn't found")
	}

	return proxy.Do(c, location)
}
