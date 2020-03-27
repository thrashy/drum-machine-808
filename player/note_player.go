package player

// NotePlayer is an interface for playing notes. This allows for swapping out different types of players.
// Return an interface to allow for unit testing output.
type NotePlayer interface {
	PlayNotes(noteNames []string) interface{}
}
