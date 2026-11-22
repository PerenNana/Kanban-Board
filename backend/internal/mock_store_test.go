package internal

import (
	"context"
	"errors"
	"sort"
)

var errNotFound = errors.New("not found")

// mockStore hält alle Daten einfach im Arbeitsspeicher (Maps statt Tabellen).
type mockStore struct {
	boards     map[int]Board
	columns    map[int]Column
	cards      map[int]Card
	nextCardID int

	// failOn: Methodenname → Fehler, den diese Methode zurückgeben soll.
	// So können Tests gezielt Datenbankfehler simulieren.
	failOn map[string]error
}

// fail liefert den eingetragenen Fehler für eine Methode (oder nil).
func (m *mockStore) fail(method string) error {
	return m.failOn[method]
}

// Jeder Test bekommt einen frischen Mock mit denselben Startdaten.
func newMockStore() *mockStore {
	return &mockStore{
		boards: map[int]Board{
			1: {ID: 1, Name: "Testboard"},
		},
		columns: map[int]Column{
			1: {ID: 1, BoardID: 1, Name: "To Do", Position: 1},
			2: {ID: 2, BoardID: 1, Name: "Done", Position: 2},
		},
		cards: map[int]Card{
			1: {ID: 1, ColumnID: 1, Title: "Erste Karte", Description: "Beschreibung", Position: 1},
		},
		nextCardID: 2,
	}
}

func (m *mockStore) ListBoards(ctx context.Context) ([]Board, error) {
	if err := m.fail("ListBoards"); err != nil {
		return nil, err
	}
	var out []Board
	for _, b := range m.boards {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *mockStore) GetBoardByID(ctx context.Context, id int) (Board, error) {
	if err := m.fail("GetBoardByID"); err != nil {
		return Board{}, err
	}
	b, ok := m.boards[id]
	if !ok {
		return Board{}, errNotFound
	}
	return b, nil
}

func (m *mockStore) GetColumnByID(ctx context.Context, id int) (Column, error) {
	if err := m.fail("GetColumnByID"); err != nil {
		return Column{}, err
	}
	c, ok := m.columns[id]
	if !ok {
		return Column{}, errNotFound
	}
	return c, nil
}

func (m *mockStore) GetColumnsByBoardID(ctx context.Context, boardID int) ([]Column, error) {
	if err := m.fail("GetColumnsByBoardID"); err != nil {
		return nil, err
	}
	var out []Column
	for _, c := range m.columns {
		if c.BoardID == boardID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Position < out[j].Position })
	return out, nil
}

func (m *mockStore) GetCardByID(ctx context.Context, id int) (Card, error) {
	if err := m.fail("GetCardByID"); err != nil {
		return Card{}, err
	}
	c, ok := m.cards[id]
	if !ok {
		return Card{}, errNotFound
	}
	return c, nil
}

func (m *mockStore) GetCardsByColumnID(ctx context.Context, columnID int) ([]Card, error) {
	if err := m.fail("GetCardsByColumnID"); err != nil {
		return nil, err
	}
	var out []Card
	for _, c := range m.cards {
		if c.ColumnID == columnID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Position < out[j].Position })
	return out, nil
}

func (m *mockStore) CreateCard(ctx context.Context, c Card) (Card, error) {
	if err := m.fail("CreateCard"); err != nil {
		return Card{}, err
	}
	c.ID = m.nextCardID
	m.nextCardID++
	m.cards[c.ID] = c
	return c, nil
}

func (m *mockStore) UpdateCard(ctx context.Context, c Card) error {
	if err := m.fail("UpdateCard"); err != nil {
		return err
	}
	if _, ok := m.cards[c.ID]; !ok {
		return errNotFound
	}
	m.cards[c.ID] = c
	return nil
}
