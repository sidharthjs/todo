package notestore

import (
	"context"
)

//Note is the model for the Notes
type Note struct {
	ID        string
	Title     string
	Body      string
	UserID    string
	CreatedAt string
}

//NoteStore is the interface for the note storage
type NoteStore interface {
	Create(ctx context.Context, note Note) error
	Read(ctx context.Context, noteID, userID string) (Note, error)
	ReadAll(ctx context.Context, userID string) ([]Note, error)
	Update(ctx context.Context, note Note) error
	Delete(ctx context.Context, noteID, userID string) error
}
