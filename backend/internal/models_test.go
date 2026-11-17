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
