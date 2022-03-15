package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	migrate "local/sidharthjs/todo/db"
	"local/sidharthjs/todo/notestore"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

var (
	testDB   = DB{}
	user     = "user_1"
	password = "secret"
	database = "testDB"
)

func TestCreateNote(t *testing.T) {
	noteID1 := uuid.New().String()
	noteID2 := uuid.New().String()
	var testCases = []struct {
		description   string
		inputNote     notestore.Note
		expectedNote  notestore.Note
		expectedError error
	}{
		{
			description: "Create note test case 1",
			inputNote: notestore.Note{
				ID:     noteID1,
				Title:  "Note 1",
				Body:   "This is a sample Note 1",
				UserID: "user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID1,
				Title:  "Note 1",
				Body:   "This is a sample Note 1",
				UserID: "user_1",
			},
			expectedError: nil,
		},
		{
			description: "Create note test case 2",
			inputNote: notestore.Note{
				ID:     noteID2,
				Title:  "Note 2",
				Body:   "This is a sample Note 2",
				UserID: "user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID2,
				Title:  "Note 2",
				Body:   "This is a sample Note 2",
				UserID: "user_1",
			},
			expectedError: nil,
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		err := testDB.Create(context.Background(), testCase.inputNote)
		assert.Nil(err)

		note, err := testDB.Read(context.Background(), testCase.inputNote.ID, testCase.inputNote.UserID)
		assert.Nil(err)
		assert.Equal(testCase.expectedNote.ID, note.ID)
		assert.Equal(testCase.expectedNote.Title, note.Title)
		assert.Equal(testCase.expectedNote.Body, note.Body)
		assert.Equal(testCase.expectedNote.UserID, note.UserID)
	}
}

func TestReadNote(t *testing.T) {
	noteID1 := uuid.New().String()
	noteID2 := uuid.New().String()
	var testCases = []struct {
		description   string
		inputNote     notestore.Note
		expectedNote  notestore.Note
		expectedError error
	}{
		{
			description: "Read note test case 1",
			inputNote: notestore.Note{
				ID:     noteID1,
				Title:  "Read Note 1",
				Body:   "This is a sample Note to be tested.",
				UserID: "user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID1,
				Title:  "Read Note 1",
				Body:   "This is a sample Note to be tested.",
				UserID: "user_1",
			},
			expectedError: nil,
		},
		{
			description: "Read note test case 1",
			inputNote: notestore.Note{
				ID:     noteID2,
				Title:  "Read Note 2",
				Body:   "This is a sample Note to be tested - 2",
				UserID: "user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID2,
				Title:  "Read Note 2",
				Body:   "This is a sample Note to be tested - 2",
				UserID: "user_1",
			},
			expectedError: nil,
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		err := testDB.Create(context.Background(), testCase.inputNote)
		assert.Nil(err)

		note, err := testDB.Read(context.Background(), testCase.inputNote.ID, testCase.inputNote.UserID)
		assert.Nil(err)
		assert.Equal(testCase.expectedNote.ID, note.ID)
		assert.Equal(testCase.expectedNote.Title, note.Title)
		assert.Equal(testCase.expectedNote.Body, note.Body)
		assert.Equal(testCase.expectedNote.UserID, note.UserID)
	}
}

func TestReadNotes(t *testing.T) {
	noteID1 := uuid.New().String()
	noteID2 := uuid.New().String()
	noteID3 := uuid.New().String()
	noteID4 := uuid.New().String()
	userID1 := "ReadAll_user_1"
	userID2 := "ReadAll_user_2"
	var testCases = []struct {
		description   string
		inputNotes    []notestore.Note
		expectedNotes []notestore.Note
		expectedError error
	}{
		{
			description: "Read notes test case 1 for user 1",
			inputNotes: []notestore.Note{
				{
					ID:     noteID1,
					Title:  "Sample note 1",
					Body:   "This is a sample Note to be tested - 1",
					UserID: userID1,
				},
				{
					ID:     noteID2,
					Title:  "Sample note 2",
					Body:   "This is a sample Note to be tested - 2",
					UserID: userID1,
				},
			},
			expectedNotes: []notestore.Note{
				{
					ID:     noteID1,
					Title:  "Sample note 1",
					Body:   "This is a sample Note to be tested - 1",
					UserID: userID1,
				},
				{
					ID:     noteID2,
					Title:  "Sample note 2",
					Body:   "This is a sample Note to be tested - 2",
					UserID: userID1,
				},
			},
			expectedError: nil,
		},
		{
			description: "Read notes test case 2 for user 2",
			inputNotes: []notestore.Note{
				{
					ID:     noteID3,
					Title:  "Sample note 3",
					Body:   "This is a sample Note to be tested - 3",
					UserID: userID2,
				},
				{
					ID:     noteID4,
					Title:  "Sample note 4",
					Body:   "This is a sample Note to be tested - 4",
					UserID: userID2,
				},
			},
			expectedNotes: []notestore.Note{
				{
					ID:     noteID3,
					Title:  "Sample note 3",
					Body:   "This is a sample Note to be tested - 3",
					UserID: userID2,
				},
				{
					ID:     noteID4,
					Title:  "Sample note 4",
					Body:   "This is a sample Note to be tested - 4",
					UserID: userID2,
				},
			},
			expectedError: nil,
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		for _, inputNote := range testCase.inputNotes {
			err := testDB.Create(context.Background(), inputNote)
			assert.Nil(err)
		}

		notes, err := testDB.ReadAll(context.Background(), testCase.inputNotes[0].UserID)
		assert.Nil(err)
		for i := 0; i < len(testCase.expectedNotes); i++ {
			assert.Equal(testCase.expectedNotes[i].ID, notes[i].ID)
			assert.Equal(testCase.expectedNotes[i].Title, notes[i].Title)
			assert.Equal(testCase.expectedNotes[i].Body, notes[i].Body)
			assert.Equal(testCase.expectedNotes[i].UserID, notes[i].UserID)
		}
	}
}

func TestUpdateNote(t *testing.T) {
	noteID1 := uuid.New().String()
	noteID2 := uuid.New().String()
	var testCases = []struct {
		description   string
		initialNote   notestore.Note
		updateNote    notestore.Note
		expectedNote  notestore.Note
		expectedError error
	}{
		{
			description: "Update note test case 1",
			initialNote: notestore.Note{
				ID:     noteID1,
				Title:  "Update Note 1",
				Body:   "This is a sample Note to be updated.",
				UserID: "Update_user_1",
			},
			updateNote: notestore.Note{
				ID:     noteID1,
				Title:  "Updated Note 1",
				Body:   "This note is updated",
				UserID: "Update_user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID1,
				Title:  "Updated Note 1",
				Body:   "This note is updated",
				UserID: "Update_user_1",
			},
			expectedError: nil,
		},
		{
			description: "Update note test case 2",
			initialNote: notestore.Note{
				ID:     noteID2,
				Title:  "My favourite note",
				Body:   "My favourites",
				UserID: "Update_user_1",
			},
			updateNote: notestore.Note{
				ID:     noteID2,
				Title:  "My favourite note - updated",
				Body:   "My favourite are updated now.",
				UserID: "Update_user_1",
			},
			expectedNote: notestore.Note{
				ID:     noteID2,
				Title:  "My favourite note - updated",
				Body:   "My favourite are updated now.",
				UserID: "Update_user_1",
			},
			expectedError: nil,
		},
		{
			description: "Don't update when user does not own the note",
			initialNote: notestore.Note{
				ID:     noteID2,
				Title:  "Update_user_2's note",
				Body:   "This note belongs to Update_user_2",
				UserID: "Update_user_2",
			},
			updateNote: notestore.Note{
				ID:     noteID2,
				Title:  "My favourite note - updated",
				Body:   "My favourite are updated now.",
				UserID: "Update_user_3",
			},
			expectedNote: notestore.Note{
				ID:     noteID2,
				Title:  "Update_user_2's note",
				Body:   "This note belongs to Update_user_2",
				UserID: "Update_user_2",
			},
			expectedError: fmt.Errorf("rows affected for update note call is 0"),
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		err := testDB.Create(context.Background(), testCase.initialNote)
		assert.Nil(err)

		if testCase.initialNote.UserID == testCase.updateNote.UserID {
			err = testDB.Update(context.Background(), testCase.updateNote)
			assert.Nil(err)
		} else {
			err = testDB.Update(context.Background(), testCase.updateNote)
			assert.EqualError(err, testCase.expectedError.Error())
		}

		note, err := testDB.Read(context.Background(), testCase.initialNote.ID, testCase.initialNote.UserID)
		assert.Nil(err)
		assert.Equal(testCase.expectedNote.ID, note.ID)
		assert.Equal(testCase.expectedNote.Title, note.Title)
		assert.Equal(testCase.expectedNote.Body, note.Body)
		assert.Equal(testCase.expectedNote.UserID, note.UserID)
	}
}

func TestDeleteNote(t *testing.T) {
	noteID1 := uuid.New().String()
	noteID2 := uuid.New().String()
	var testCases = []struct {
		description   string
		inputNote     notestore.Note
		expectedError error
	}{
		{
			description: "Delete note test case 1",
			inputNote: notestore.Note{
				ID:     noteID1,
				Title:  "Delete Note 1",
				Body:   "This is a sample Note to be deleted.",
				UserID: "user_1",
			},
			expectedError: fmt.Errorf("Note '%s' is not found", noteID1),
		},
		{
			description: "Delete note test case 2",
			inputNote: notestore.Note{
				ID:     noteID2,
				Title:  "Delete Note 2",
				Body:   "This is a sample Note to be deleted.",
				UserID: "user_1",
			},
			expectedError: fmt.Errorf("Note '%s' is not found", noteID2),
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		err := testDB.Create(context.Background(), testCase.inputNote)
		assert.Nil(err)

		err = testDB.Delete(context.Background(), testCase.inputNote.ID, testCase.inputNote.UserID)
		assert.Nil(err)

		_, err = testDB.Read(context.Background(), testCase.inputNote.ID, testCase.inputNote.UserID)
		assert.EqualError(err, testCase.expectedError.Error())
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + user,
			"POSTGRES_DB=" + database,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, hostAndPort, database)
	fmt.Println("hostAndPort:", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		fmt.Println("connecting...")
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		testDB = DB{db}
		err = db.Ping()
		if err != nil {
			fmt.Println("unable to ping:", err)
			return err
		}
		fmt.Println("Pinged db successfully!")
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = migrate.Migrate(databaseUrl, "file://../../db/migrations")
	if err != nil {
		log.Fatalf("error performing db migration: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
