package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"local/sidharthjs/todo/notestore"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB struct that represents the Client
type DB struct {
	*sql.DB
}

// NewClient returns the DB client and error
func NewClient(dbURL string) (*DB, error) {
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %s", err)
	}

	return &DB{conn}, nil
}

//Create creates a note in the DB
func (db *DB) Create(ctx context.Context, note notestore.Note) error {
	sql := "INSERT INTO notes(user_id, id, title, body, created_at) VALUES($1, $2, $3, $4, $5);"
	ct, err := db.ExecContext(ctx, sql, note.UserID, note.ID, note.Title, note.Body, time.Now())
	if err != nil {
		return fmt.Errorf("unable to store note '%s': %s", note.ID, err)
	}

	n, err := ct.RowsAffected()
	if err != nil {
		return fmt.Errorf("error in getting rows affected: %s", err)
	}
	if n == 0 {
		return fmt.Errorf("rows affected for store note call is 0")
	}
	return nil
}

//Read reads a note from the DB
func (db *DB) Read(ctx context.Context, noteID, userID string) (notestore.Note, error) {
	sqlQuery := "SELECT id, title, body, user_id, created_at FROM notes WHERE id=$1 and user_id=$2;"
	row := db.QueryRowContext(ctx, sqlQuery, noteID, userID)

	var note notestore.Note
	err := row.Scan(&note.ID, &note.Title, &note.Body, &note.UserID, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return notestore.Note{}, fmt.Errorf("Note '%s' is not found", noteID)
		}
		return notestore.Note{}, fmt.Errorf("error occurred while retrieving the note: %s", err)
	}

	return note, nil
}

//ReadAll reads all the notes for the given user ID
func (db *DB) ReadAll(ctx context.Context, userID string) ([]notestore.Note, error) {
	sql := "SELECT id, title, body, user_id, created_at FROM notes WHERE user_id=$1;"
	rows, err := db.QueryContext(ctx, sql, userID)
	if err != nil {
		return []notestore.Note{}, fmt.Errorf("error occurred while querying the note: %s", err)
	}
	defer rows.Close()

	var notes []notestore.Note
	for rows.Next() {
		var note notestore.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.UserID, &note.CreatedAt)
		if err != nil {
			return []notestore.Note{}, fmt.Errorf("error occurred while scanning the rows: %s", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

//Update updates a note
func (db *DB) Update(ctx context.Context, note notestore.Note) error {
	sql := "UPDATE notes SET title=$1, body=$2 WHERE id=$3 AND user_id=$4;"
	ct, err := db.ExecContext(ctx, sql, note.Title, note.Body, note.ID, note.UserID)
	if err != nil {
		return fmt.Errorf("unable to update note '%s': %s", note.ID, err)
	}

	n, err := ct.RowsAffected()
	if err != nil {
		return fmt.Errorf("error in getting rows affected: %s", err)
	}
	if n == 0 {
		return fmt.Errorf("rows affected for update note call is 0")
	}
	return nil
}

// Delete deletes a note
func (db *DB) Delete(ctx context.Context, noteID, userID string) error {
	sql := "DELETE FROM notes WHERE id=$1 AND user_id=$2;"
	ct, err := db.ExecContext(ctx, sql, noteID, userID)
	if err != nil {
		return fmt.Errorf("unable to delete note '%s': %s", noteID, err)
	}

	n, err := ct.RowsAffected()
	if err != nil {
		return fmt.Errorf("error in getting rows affected: %s", err)
	}
	if n == 0 {
		return fmt.Errorf("rows affected for delete note call is 0")
	}
	return nil
}
