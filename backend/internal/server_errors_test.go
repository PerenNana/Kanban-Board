package internal

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// errDB simuliert einen beliebigen Datenbankfehler.
var errDB = errors.New("db kaputt")

// --- NewServer: ungültiges Template-Muster ---

func TestNewServerInvalidTemplates(t *testing.T) {
	// Ein Muster, das auf keine Datei passt, führt bei ParseGlob zu einem Fehler
	_, err := NewServer(newMockStore(), "gibt-es-nicht/*.html")
	if err == nil {
		t.Fatal("Fehler erwartet, aber nil bekommen")
	}
}

// --- render: Template existiert nicht ---

func TestRenderTemplateError(t *testing.T) {
	// Server ohne die benötigten Templates → ExecuteTemplate schlägt fehl
	srv := &Server{
		db:   newMockStore(),
		tmpl: template.Must(template.New("leer").Parse("")),
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- GET-Seiten: Fehlerzweige ---

func TestGetPagesErrors(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		failOn     map[string]error
		wantStatus int
	}{
		{"Startseite: ListBoards schlägt fehl", "/",
			map[string]error{"ListBoards": errDB}, http.StatusInternalServerError},
		{"Board: GetColumnsByBoardID schlägt fehl", "/boards/1",
			map[string]error{"GetColumnsByBoardID": errDB}, http.StatusInternalServerError},
		{"Board: GetCardsByColumnID schlägt fehl", "/boards/1",
			map[string]error{"GetCardsByColumnID": errDB}, http.StatusInternalServerError},
		{"Bearbeiten: ungültige Karten-ID", "/cards/abc/edit",
			nil, http.StatusNotFound},
		{"Bearbeiten: GetColumnByID schlägt fehl", "/cards/1/edit",
			map[string]error{"GetColumnByID": errDB}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, store := newTestServer(t)
			store.failOn = tt.failOn

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			srv.Routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// --- POST /cards/{id}: Karte aktualisieren – Fehlerzweige ---

func TestUpdateCardErrors(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		form       url.Values
		failOn     map[string]error
		wantStatus int
	}{
		{"Ungültige Karten-ID", "/cards/abc",
			url.Values{"title": {"Titel"}}, nil, http.StatusNotFound},
		{"Leerer Titel", "/cards/1",
			url.Values{"title": {"   "}}, nil, http.StatusBadRequest},
		{"Unbekannte Karte", "/cards/99",
			url.Values{"title": {"Titel"}}, nil, http.StatusNotFound},
		{"UpdateCard schlägt fehl", "/cards/1",
			url.Values{"title": {"Titel"}},
			map[string]error{"UpdateCard": errDB}, http.StatusInternalServerError},
		{"GetColumnByID schlägt fehl", "/cards/1",
			url.Values{"title": {"Titel"}},
			map[string]error{"GetColumnByID": errDB}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, store := newTestServer(t)
			store.failOn = tt.failOn
			before := store.cards[1]

			rec := postForm(srv, tt.path, tt.form)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", rec.Code, tt.wantStatus)
			}
			// Bei Validierungsfehlern darf die Karte nicht verändert werden
			if tt.wantStatus != http.StatusInternalServerError && store.cards[1] != before {
				t.Errorf("Karte wurde trotz Fehler verändert: %+v", store.cards[1])
			}
		})
	}
}

// --- POST /cards: Karte erstellen – Fehlerzweige ---

func TestCreateCardErrors(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		failOn     map[string]error
		wantStatus int
	}{
		{"Spalte fehlt",
			url.Values{"title": {"Titel"}}, nil, http.StatusBadRequest},
		{"Ungültige Spalten-ID",
			url.Values{"column_id": {"abc"}, "title": {"Titel"}}, nil, http.StatusBadRequest},
		{"Unbekannte Spalte",
			url.Values{"column_id": {"99"}, "title": {"Titel"}}, nil, http.StatusNotFound},
		{"GetCardsByColumnID schlägt fehl",
			url.Values{"column_id": {"1"}, "title": {"Titel"}},
			map[string]error{"GetCardsByColumnID": errDB}, http.StatusInternalServerError},
		{"CreateCard schlägt fehl",
			url.Values{"column_id": {"1"}, "title": {"Titel"}},
			map[string]error{"CreateCard": errDB}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, store := newTestServer(t)
			store.failOn = tt.failOn

			rec := postForm(srv, "/cards", tt.form)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if len(store.cards) != 1 {
				t.Errorf("es wurde trotzdem eine Karte angelegt: %d Karten", len(store.cards))
			}
		})
	}
}
