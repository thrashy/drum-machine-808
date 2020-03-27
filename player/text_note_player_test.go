package player

import (
	"fmt"
	"testing"
)

// TestPlayNotes tests the functionality of TextNotePlayer.PlayNotes()
func TestTextNotePlayerPlayNotes(t *testing.T) {

	const beamedEightNote = string(9835) // ♫

	player := NewTextNotePlayer()

	testNotes := []string{
		"hh",
		"unknown",
	}

	actual := player.PlayNotes(testNotes)
	expected := fmt.Sprintf("%s %s %s", beamedEightNote, "Closed Hi-Hat, unknown", beamedEightNote)

	if actual != expected {
		t.Errorf("Notes played text did not match, actual: %s, expected: %s.", actual, expected)
	}
}
