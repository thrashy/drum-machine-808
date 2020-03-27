package player

import (
	"testing"
)

// TestPlayNotes tests the functionality of WavNotePlayer.PlayNotes()
func TestWavNotePlayerPlayNotes(t *testing.T) {

	// Simple test to make sure we don't panic(). No output is returned for this call. Only audio is played.
	player := NewWavNotePlayer()
	player.PlayNotes([]string{
		"sn",
		"bd",
		"hh",
	})
}
