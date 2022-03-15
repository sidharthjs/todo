package main

import (
	"os"

	migrate "local/sidharthjs/todo/db"
	"local/sidharthjs/todo/handlers/authhandler"
	"local/sidharthjs/todo/handlers/noteshandler"
	"local/sidharthjs/todo/middleware"
	"local/sidharthjs/todo/notestore/postgres"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Set log level to debug
	log.SetLevel(log.DebugLevel)

	// Perform DB Migration
	postgresUser := readEnv("POSTGRES_USER")
	postgresPassword := readEnv("POSTGRES_PASSWORD")
	postgresDB := readEnv("POSTGRES_DB")
	postgresHost := readEnv("POSTGRES_HOST")
	postgresPort := readEnv("POSTGRES_PORT")
	dbURL := "postgres://" + postgresUser + ":" + postgresPassword + "@" + postgresHost + ":" + postgresPort + "/" + postgresDB + "?sslmode=disable"
	err := migrate.Migrate(dbURL, "file://db/migrations")
	if err != nil {
		log.Fatalf("error performing db migration: %s", err)
	}

	// Init DB
	db, err := postgres.NewClient(dbURL)
	if err != nil {
		log.Fatalf("error creating db client: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("unable to ping the db")
	}

	// Init handlers
	authHandler := authhandler.New(readEnv("GITHUB_CLIENT_ID"), readEnv("GITHUB_CLIENT_SECRET"),
		readEnv("LOGIN_URI"), readEnv("ACCESS_TOKEN_URI"),
		readEnv("REDIRECT_URI"), readEnv("PROFILE_URI"), readEnv("USERS_SVC_ENDPOINT"))
	notesHandler := noteshandler.New(db)

	// Define routes
	app := fiber.New()
	app.Static("/", "./public/login.html")
	app.Get("/login/github", authHandler.InitiateOAuth)
	app.Get("/github/callback", authHandler.ProcessCallback)

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"POSTGRES_USER":        readEnv("POSTGRES_USER"),
			"POSTGRES_PASSWORD":    readEnv("POSTGRES_PASSWORD"),
			"POSTGRES_DB":          readEnv("POSTGRES_DB"),
			"POSTGRES_HOST":        readEnv("POSTGRES_HOST"),
			"POSTGRES_PORT":        readEnv("POSTGRES_PORT"),
			"GITHUB_CLIENT_ID":     readEnv("GITHUB_CLIENT_ID"),
			"GITHUB_CLIENT_SECRET": readEnv("GITHUB_CLIENT_SECRET"),
			"LOGIN_URI":            readEnv("LOGIN_URI"),
			"ACCESS_TOKEN_URI":     readEnv("ACCESS_TOKEN_URI"),
			"REDIRECT_URI":         readEnv("REDIRECT_URI"),
			"PROFILE_URI":          readEnv("PROFILE_URI"),
			"USERS_SVC_ENDPOINT":   readEnv("USERS_SVC_ENDPOINT"),
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.Redirect(readEnv("USERS_SVC_ENDPOINT") + "/users")
	})

	middleware.SetupAuthentication(app)

	app.Get("/notes/:note_id", notesHandler.ReadNote)
	app.Put("/notes/:note_id", notesHandler.UpdateNote)
	app.Get("/notes", notesHandler.ReadNotes)
	app.Post("/notes", notesHandler.CreateNote)
	app.Delete("/notes/:note_id", notesHandler.DeleteNote)

	log.Info("app running...")
	log.Fatal(app.Listen(":4010"))
}

func readEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("env variable %s is not set or empty", key)
	}
	return val
}
