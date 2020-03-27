package player

import (
	"fmt"
	"strings"
)

// TextNotePlayer defines a struct for playing notes using a text player.
type TextNotePlayer struct {
	NoteMap map[string]string
}

// NewTextNotePlayer initializes a TextNotePlayer with a mapping for note abbreviations to their full names.
func NewTextNotePlayer() *TextNotePlayer {

	return &TextNotePlayer{
		NoteMap: map[string]string{
			"bd": "Bass Drum",
			"sn": "Snare Drum",
			"hh": "Closed Hi-Hat",
			"oh": "Open Hi-Hat",
			"rd": "Ride Cymbal",
		},
	}
}

// PlayNotes uses standard output to print the names of notes being played.
func (player *TextNotePlayer) PlayNotes(noteNames []string) interface{} {

	const concatString = ", "
	const beamedEightNote = string(9835) // ♫

	output := ""

	for _, note := range noteNames {
		noteFullName, exists := player.NoteMap[note]

		// if the abbreviation does not exist in the mapping, default to the abbreviation name
		if !exists {
			noteFullName = note
		}

		// use naive string concatenation method, as the number of notes in the slice is small
		output += noteFullName + concatString
	}

	output = fmt.Sprintf("%s %s %s", beamedEightNote, strings.TrimSuffix(output, concatString), beamedEightNote)
	fmt.Println(output)

	return output
}
