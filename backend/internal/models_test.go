package internal

import "testing"

func TestBoardString(t *testing.T) {
	// Vorbereiten: Ein Board-Objekt erstellen
	b := Board{ID: 1, Name: "Alpha"}

	// Ausführen: String-Methode aufrufen
	got := b.String()
	want := `Board{ID:1, Name:"Alpha"}`

	// Prüfen: Ergebnis vergleichen
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestColumnString(t *testing.T) {
	// Vorbereiten: Ein Column-Objekt erstellen
	c := Column{ID: 2, BoardID: 1, Name: "To Do", Position: 1}

	// Ausführen: String-Methode aufrufen
	got := c.String()
	want := `Column{ID:2, BoardID:1, Name:"To Do", Position:1}`

	// Prüfen: Ergebnis vergleichen
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestCardString(t *testing.T) {
	// Vorbereiten: Ein Card-Objekt erstellen
	c := Card{ID: 3, ColumnID: 2, Title: "Task", Description: "Nevahh! Do such a thing", Position: 2}

	// Ausführen: String-Methode aufrufen
	got := c.String()
	want := `Card{ID:3, ColumnID:2, Title:"Task", Description:"Nevahh! Do such a thing", Position:2}`

	// Prüfen: Ergebnis vergleichen
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
