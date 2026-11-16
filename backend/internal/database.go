package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) CreateCard(ctx context.Context, c Card) (Card, error) {
	var id int
	err := db.pool.QueryRow(ctx,
		`INSERT INTO cards (column_id, title, description, position)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		c.ColumnID, c.Title, c.Description, c.Position,
	).Scan(&id)
	if err != nil {
		return Card{}, fmt.Errorf("CreateCard: %w", err)
	}
	c.ID = id
	return c, nil
}

func (db *DB) GetCardByID(ctx context.Context, id int) (Card, error) {
	var c Card
	err := db.pool.QueryRow(ctx,
		`SELECT id, column_id, title, description, position
		 FROM cards WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Card{}, fmt.Errorf("GetCardByID: not found")
		}
		return Card{}, fmt.Errorf("GetCardByID: %w", err)
	}
	return c, nil
}

func (db *DB) UpdateCard(ctx context.Context, c Card) error {
	ct, err := db.pool.Exec(ctx,
		`UPDATE cards SET column_id = $1, title = $2, description = $3, position = $4
		 WHERE id = $5`,
		c.ColumnID, c.Title, c.Description, c.Position, c.ID)
	if err != nil {
		return fmt.Errorf("UpdateCard: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("UpdateCard: not found")
	}
	return nil
}

func (db *DB) DeleteCard(ctx context.Context, id int) error {
	ct, err := db.pool.Exec(ctx, `DELETE FROM cards WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("DeleteCard: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("DeleteCard: not found")
	}
	return nil
}

func (db *DB) CreateBoard(ctx context.Context, name string) (Board, error) {
	var b Board
	err := db.pool.QueryRow(ctx,
		`INSERT INTO boards (name) VALUES ($1) RETURNING id, name`,
		name,
	).Scan(&b.ID, &b.Name)
	if err != nil {
		return Board{}, fmt.Errorf("CreateBoard: %w", err)
	}
	return b, nil
}

func (db *DB) CreateColumn(ctx context.Context, col Column) (Column, error) {
	err := db.pool.QueryRow(ctx,
		`INSERT INTO columns (board_id, name, position) VALUES ($1, $2, $3) RETURNING id`,
		col.BoardID, col.Name, col.Position,
	).Scan(&col.ID)
	if err != nil {
		return Column{}, fmt.Errorf("CreateColumn: %w", err)
	}
	return col, nil
}

func (db *DB) ListBoards(ctx context.Context) ([]Board, error) {
	rows, err := db.pool.Query(ctx, `SELECT id, name FROM boards ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("ListBoards: %w", err)
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, fmt.Errorf("ListBoards scan: %w", err)
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

func (db *DB) GetColumnsByBoardID(ctx context.Context, boardID int) ([]Column, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, board_id, name, position FROM columns
		 WHERE board_id = $1 ORDER BY position`, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetColumnsByBoardID: %w", err)
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var c Column
		if err := rows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position); err != nil {
			return nil, fmt.Errorf("GetColumnsByBoardID scan: %w", err)
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

func (db *DB) GetCardsByColumnID(ctx context.Context, columnID int) ([]Card, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, column_id, title, description, position FROM cards
		 WHERE column_id = $1 ORDER BY position`, columnID)
	if err != nil {
		return nil, fmt.Errorf("GetCardsByColumnID: %w", err)
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position); err != nil {
			return nil, fmt.Errorf("GetCardsByColumnID scan: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (db *DB) GetBoardByID(ctx context.Context, id int) (Board, error) {
	var b Board
	err := db.pool.QueryRow(ctx,
		`SELECT id, name FROM boards WHERE id = $1`, id,
	).Scan(&b.ID, &b.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Board{}, fmt.Errorf("GetBoardByID: not found")
		}
		return Board{}, fmt.Errorf("GetBoardByID: %w", err)
	}
	return b, nil
}

func (db *DB) GetColumnByID(ctx context.Context, id int) (Column, error) {
	var c Column
	err := db.pool.QueryRow(ctx,
		`SELECT id, board_id, name, position FROM columns WHERE id = $1`, id,
	).Scan(&c.ID, &c.BoardID, &c.Name, &c.Position)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Column{}, fmt.Errorf("GetColumnByID: not found")
		}
		return Column{}, fmt.Errorf("GetColumnByID: %w", err)
	}
	return c, nil
}
