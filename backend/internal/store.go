package internal

import "context"

// Store beschreibt alles, was der Server von einer Datenbank braucht.
type Store interface {
	ListBoards(ctx context.Context) ([]Board, error)
	GetBoardByID(ctx context.Context, id int) (Board, error)
	GetColumnByID(ctx context.Context, id int) (Column, error)
	GetColumnsByBoardID(ctx context.Context, boardID int) ([]Column, error)
	GetCardByID(ctx context.Context, id int) (Card, error)
	GetCardsByColumnID(ctx context.Context, columnID int) ([]Card, error)
	CreateCard(ctx context.Context, c Card) (Card, error)
	UpdateCard(ctx context.Context, c Card) error
}

// Prüft beim Kompilieren, dass *DB wirklich alle Methoden hat.
var _ Store = (*DB)(nil)
