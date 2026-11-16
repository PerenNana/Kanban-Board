package internal

import "fmt"

type Board struct {
	ID   int
	Name string
}

type Column struct {
	ID       int
	BoardID  int
	Name     string
	Position int
}

type Card struct {
	ID          int
	ColumnID    int
	Title       string
	Description string
	Position    int
}

// Daten für board.html
type ColumnView struct {
	Column Column
	Cards  []Card
}

type BoardView struct {
	Board   Board
	Columns []ColumnView
}

// Daten für card_edit.html
type CardEditView struct {
	Card    Card
	BoardID int
}

// --- Wie soll ein Objekt als Text aussehen? ---

func (b Board) String() string {
	return fmt.Sprintf("Board{ID:%d, Name:%q}", b.ID, b.Name)
}

func (c Column) String() string {
	return fmt.Sprintf("Column{ID:%d, BoardID:%d, Name:%q, Position:%d}",
		c.ID, c.BoardID, c.Name, c.Position)
}

func (c Card) String() string {
	return fmt.Sprintf("Card{ID:%d, ColumnID:%d, Title:%q, Description:%q, Position:%d}",
		c.ID, c.ColumnID, c.Title, c.Description, c.Position)
}
