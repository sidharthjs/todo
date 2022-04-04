package noteshandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwtutil "local/sidharthjs/todo/jwt"
	"local/sidharthjs/todo/notestore"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//NotesHandler struct defintion
type NotesHandler struct {
	Store        notestore.NoteStore
	UsersService string
}

//New returns NotesHandler
func New(notestore notestore.NoteStore, usersService string) *NotesHandler {
	return &NotesHandler{
		Store:        notestore,
		UsersService: usersService,
	}
}

//CreateNote is the handler method for creating a note
func (nh *NotesHandler) CreateNote(c *fiber.Ctx) error {
	userID, _, err := jwtutil.GetUserFromJWTToken(c.Locals("user").(*jwt.Token))
	if err != nil {
		log.Errorf("error in reading user details in jwt token: %s", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	type request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	var req request
	err = c.BodyParser(&req)
	if err != nil {
		log.Errorf("unable to parse the request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to create a note",
		})
	}

	note := notestore.Note{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Body:   req.Body,
		UserID: userID,
	}

	err = nh.Store.Create(c.UserContext(), note)
	if err != nil {
		log.Errorf("unable to create a note: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to create a note",
		})
	}

	log.Infof("note %s created successfully", note.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": fmt.Sprintf("note '%s' created successfully", note.ID),
	})
}

//ReadNote is the handler method for reading a note
func (nh *NotesHandler) ReadNote(c *fiber.Ctx) error {

	userID, _, err := jwtutil.GetUserFromJWTToken(c.Locals("user").(*jwt.Token))
	if err != nil {
		log.Errorf("error in reading user details in jwt token: %s", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	noteID := c.Params("note_id")
	note, err := nh.Store.Read(c.UserContext(), noteID, userID)
	if err != nil {
		log.Errorf("unable to read note '%s': %s", noteID, err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error in reading the note",
		})
	}

	resp, err := http.Get(nh.UsersService + "/users/" + userID)
	if err != nil {
		log.Errorf("error while reading user: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error while reading user",
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error while reading user: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error while reading user",
		})
	}

	type userResponse struct {
		ID        string
		Username  string
		CreatedAt string
	}

	var user userResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Errorf("error while unmarshalling user: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error while unmarshalling user",
		})
	}
	log.Info(user)

	type readNoteResponse struct {
		NoteID    string
		NoteTitle string
		NoteBody  string
		CreatedAt string
		UserID    string
		Username  string
	}

	res := readNoteResponse{
		NoteID:    note.ID,
		NoteTitle: note.Title,
		NoteBody:  note.Body,
		CreatedAt: note.CreatedAt,
		UserID:    user.ID,
		Username:  user.Username,
	}
	log.Info(res)

	return c.Status(fiber.StatusOK).JSON(res)
}

//ReadNotes is the handler method for reading multiple notes
func (nh *NotesHandler) ReadNotes(c *fiber.Ctx) error {

	userID, _, err := jwtutil.GetUserFromJWTToken(c.Locals("user").(*jwt.Token))
	if err != nil {
		log.Errorf("error in reading user details in jwt token: %s", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	notes, err := nh.Store.ReadAll(c.UserContext(), userID)
	if err != nil {
		log.Errorf("error in reading all notes: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error in reading the notes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(notes)
}

//UpdateNote is the handler method for updating a note
func (nh *NotesHandler) UpdateNote(c *fiber.Ctx) error {

	userID, _, err := jwtutil.GetUserFromJWTToken(c.Locals("user").(*jwt.Token))
	if err != nil {
		log.Errorf("error in reading user details in jwt token: %s", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	type request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	var req request
	err = c.BodyParser(&req)
	if err != nil {
		log.Errorf("unable to parse the request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to create a note",
		})
	}

	note := notestore.Note{
		Title:  req.Title,
		ID:     c.Params("note_id"),
		Body:   req.Body,
		UserID: userID,
	}

	err = nh.Store.Update(c.UserContext(), note)
	if err != nil {
		log.Errorf("unable to update note '%s': %s", note.ID, err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error in updating the note",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("note '%s' is updated successfully", note.ID),
	})
}

//DeleteNote is the handler method for deleting a note
func (nh *NotesHandler) DeleteNote(c *fiber.Ctx) error {

	userID, _, err := jwtutil.GetUserFromJWTToken(c.Locals("user").(*jwt.Token))
	if err != nil {
		log.Errorf("error in reading user details in jwt token: %s", err)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	noteID := c.Params("note_id")

	err = nh.Store.Delete(c.UserContext(), noteID, userID)
	if err != nil {
		log.Errorf("unable to delete the note '%s': %s", noteID, err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error in deleting the note",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": fmt.Sprintf("note '%s' deleted successfully", noteID),
	})
}
