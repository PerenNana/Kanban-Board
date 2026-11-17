package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// --- Hilfsfunktionen ---

func newTestServer(t *testing.T) (*Server, *mockStore) {
	t.Helper()
	store := newMockStore()
	srv, err := NewServer(store, "../templates/*.html")
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	return srv, store
}

func postForm(srv *Server, path string, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	return rec
}

// --- GET-Seiten: mehrere Fälle in einer Tabelle ---

func TestGetPages(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{"Startseite zeigt Boards", "/", http.StatusOK, "Testboard"},
		{"Board zeigt Spalten", "/boards/1", http.StatusOK, "To Do"},
		{"Board zeigt Karten", "/boards/1", http.StatusOK, "Erste Karte"},
		{"Unbekanntes Board", "/boards/99", http.StatusNotFound, ""},
		{"Ungültige Board-ID", "/boards/abc", http.StatusNotFound, ""},
		{"Bearbeiten-Formular", "/cards/1/edit", http.StatusOK, "Erste Karte"},
		{"Unbekannte Karte", "/cards/99/edit", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := newTestServer(t)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			srv.Routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("Body enthält nicht %q", tt.wantBody)
			}
		})
	}
}

// --- POST: Karte erstellen ---

func TestCreateCard(t *testing.T) {
	srv, store := newTestServer(t)

	rec := postForm(srv, "/cards", url.Values{
		"column_id":   {"1"},
		"title":       {"Neue Karte"},
		"description": {"Test"},
	})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("Status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/boards/1" {
		t.Errorf("Redirect = %q, want %q", loc, "/boards/1")
	}

	card, ok := store.cards[2]
	if !ok {
		t.Fatal("neue Karte wurde nicht gespeichert")
	}
	if card.Title != "Neue Karte" || card.ColumnID != 1 || card.Position != 2 {
		t.Errorf("gespeicherte Karte = %+v", card)
	}
}

func TestCreateCardEmptyTitle(t *testing.T) {
	srv, store := newTestServer(t)

	rec := postForm(srv, "/cards", url.Values{
		"column_id": {"1"},
		"title":     {"   "},
	})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if len(store.cards) != 1 {
		t.Errorf("es wurde trotzdem eine Karte angelegt: %d Karten", len(store.cards))
	}
}

// --- POST: Karte aktualisieren ---

func TestUpdateCard(t *testing.T) {
	srv, store := newTestServer(t)

	rec := postForm(srv, "/cards/1", url.Values{
		"title":       {"Geänderter Titel"},
		"description": {"Neue Beschreibung"},
	})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("Status = %d, want %d", rec.Code, http.StatusSeeOther)
	}

	card := store.cards[1]
	if card.Title != "Geänderter Titel" || card.Description != "Neue Beschreibung" {
		t.Errorf("Karte nicht aktualisiert: %+v", card)
	}
	if card.ColumnID != 1 || card.Position != 1 {
		t.Errorf("Spalte/Position dürfen sich nicht ändern: %+v", card)
	}
}
