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
	err := migrate.Migrate(readEnv("DB_URL"), "file://db/migrations")
	if err != nil {
		log.Fatalf("error performing db migration: %s", err)
	}

	// Init DB
	db, err := postgres.NewClient(readEnv("DB_URL"))
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
		readEnv("REDIRECT_URI"), readEnv("PROFILE_URI"))
	notesHandler := noteshandler.New(db)

	// Define routes
	app := fiber.New()
	app.Static("/", "./public/login.html")
	app.Get("/login/github", authHandler.InitiateOAuth)
	app.Get("/github/callback", authHandler.ProcessCallback)

	middleware.SetupAuthentication(app)

	app.Get("/notes/:note_id", notesHandler.ReadNote)
	app.Put("/notes/:note_id", notesHandler.UpdateNote)
	app.Get("/notes", notesHandler.ReadNotes)
	app.Post("/notes", notesHandler.CreateNote)
	app.Delete("/notes/:note_id", notesHandler.DeleteNote)

	log.Info("app running...")
	log.Fatal(app.Listen(":4000"))
}

func readEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("env variable %s is not set or empty", key)
	}
	return val
}
