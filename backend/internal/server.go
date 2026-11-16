package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	db   *DB
	tmpl *template.Template
}

func NewServer(db *DB) (*Server, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{db: db, tmpl: tmpl}, nil
}

// 🚦 Router: welche URL → welcher Handler
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("GET /boards/{id}", s.handleBoard)
	mux.HandleFunc("GET /cards/{id}/edit", s.handleCardEdit)

	mux.HandleFunc("POST /cards", s.handleCardCreate)
	mux.HandleFunc("POST /cards/{id}", s.handleCardUpdate)
	return mux
}

// 🖨️ Hilfsfunktion: Template ausführen + Fehler abfangen
func (s *Server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "Template-Fehler", http.StatusInternalServerError)
	}
}

// --- Handler ---

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	boards, err := s.db.ListBoards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, "index.html", boards)
}

func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()

	board, err := s.db.GetBoardByID(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cols, err := s.db.GetColumnsByBoardID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view := BoardView{Board: board}
	for _, c := range cols {
		cards, err := s.db.GetCardsByColumnID(ctx, c.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		view.Columns = append(view.Columns, ColumnView{Column: c, Cards: cards})
	}

	s.render(w, "board.html", view)
}

func (s *Server) handleCardEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()

	card, err := s.db.GetCardByID(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	col, err := s.db.GetColumnByID(ctx, card.ColumnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.render(w, "card_edit.html", CardEditView{Card: card, BoardID: col.BoardID})
}

func (s *Server) handleCardUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()

	// Formularwerte lesen + prüfen
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "Titel darf nicht leer sein", http.StatusBadRequest)
		return
	}
	description := strings.TrimSpace(r.FormValue("description"))

	// Aktuelle Karte aus der DB holen
	card, err := s.db.GetCardByID(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Nur das ändern, was das Formular bearbeiten darf
	card.Title = title
	card.Description = description

	if err := s.db.UpdateCard(ctx, card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Zurück zum Board
	col, err := s.db.GetColumnByID(ctx, card.ColumnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/boards/%d", col.BoardID), http.StatusSeeOther)
}

func (s *Server) handleCardCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Formularwete lesen + prüfen
	columnID, err := strconv.Atoi(r.FormValue("column_id"))
	if err != nil {
		http.Error(w, "ungültige Spalte", http.StatusBadRequest)
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "Titel darf nicht leer sein", http.StatusBadRequest)
		return
	}
	description := strings.TrimSpace(r.FormValue("description"))

	// Gibt es die Spalte überhaupt?
	col, err := s.db.GetColumnByID(ctx, columnID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Neue Karte ans Ende der Spalte setzen
	cards, err := s.db.GetCardsByColumnID(ctx, columnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.CreateCard(ctx, Card{
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Position:    len(cards) + 1,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Zurück zum Board
	http.Redirect(w, r, fmt.Sprintf("boards/%d", col.BoardID), http.StatusSeeOther)
}
