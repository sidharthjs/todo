package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	log "github.com/sirupsen/logrus"
)

const jwtSecret = "aJWTSecret"

// SetupAuthentication set authentication middleware for /notes routes
func SetupAuthentication(app *fiber.App) {
	app.Use("/notes", jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Errorf("error in middleware authentication: %s", err)
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
			return nil
		},
		SigningKey: []byte(jwtSecret),
	}))
}
