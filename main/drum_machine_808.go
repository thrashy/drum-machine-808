package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/thrashy/drum-machine-808/model"
	"github.com/thrashy/drum-machine-808/player"
)

// main is the
func main() {

	var err interface{}

	defer func() {
		if r := recover(); r != nil {
			handleApplicationError(r)
		}
	}()
	err = process()

	if err != nil {
		handleApplicationError(err)
	}
}

// process is the driver of the application. All errors should be sent back to main() for global error handling.
func process() error {

	const configLocation = "../configs/conf.toml"

	sb, err := loadConfigurationFromToml(configLocation)
	if err != nil {
		return err
	}

	if err = sb.Validate(); err != nil {
		return err
	}

	var songNames []string

	for name := range sb.SongBeatData {
		songNames = append(songNames, name)
	}

	songName, err := getValueFromPrompt(songNames, "Select Song")
	if err != nil {
		return err
	}

	outputFormats := []string{
		"Text",
		"Audio",
	}

	selectedOutput, err := getValueFromPrompt(outputFormats, "Select Output Format")
	if err != nil {
		return err
	}

	songBeat := sb.SongBeatData[songName]

	notePlayer := getPlayer(selectedOutput)

	// Add a newline before output, as promptui clears data on the next line.
	fmt.Printf("\nPlaying %s at BPM: %d\n", songName, songBeat.BeatsPerMinute)

	playSong(songBeat, notePlayer)

	fmt.Printf("Song %s completed!", songName)

	return nil
}

// getPlayer is a factory method for getting a new NotePlayer.
func getPlayer(outputFormat string) player.NotePlayer {
	switch outputFormat {
	case "Audio":
		return player.NewWavNotePlayer()
	case "Text":
		return player.NewTextNotePlayer()
	default:
		return player.NewWavNotePlayer()
	}

}

// getValueFromPrompt prompts a user for a input and then returns the selected value.
func getValueFromPrompt(items []string, label string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, songName, err := prompt.Run()

	return songName, err
}

// playSong plays a song with the defined beat pattern passed into the function.
func playSong(songBeat *model.SongBeat, player player.NotePlayer) {

	const secondsInMinute float64 = 60.0
	const songPlayTimeInSeconds = 20

	secondsPerBeat := secondsInMinute / float64(songBeat.BeatsPerMinute)
	secondsPerNote := secondsPerBeat / float64(songBeat.NotesPerBeat)

	notesToPlay := getNotesToPlay(songBeat)

	noteTicker := getNoteTicker(secondsPerNote)
	done := make(chan bool)

	go playNotes(done, noteTicker, player, notesToPlay)

	time.Sleep(songPlayTimeInSeconds * time.Second)
	noteTicker.Stop()
	done <- true
}

// getNotesToPlay takes a model value and creates a 2D array of all the notes to play. We could add it to the
// configuration like this to avoid this logic, but it is less human readable
func getNotesToPlay(songBeat *model.SongBeat) [][]string {

	notesToPlay := make([][]string, 0)
	for i := range notesToPlay {
		notesToPlay[i] = make([]string, 0)
	}

	numBeats := songBeat.BeatsPerSequence * songBeat.NotesPerBeat

	for i := 0; i < numBeats; i++ {
		soundsInBeat := make([]string, 0)
		for beatType := range songBeat.BeatPattern {
			if songBeat.BeatPattern[beatType][i] {
				soundsInBeat = append(soundsInBeat, beatType)
			}
		}
		notesToPlay = append(notesToPlay, soundsInBeat)
	}

	return notesToPlay
}

// getNoteTicker returns a new ticker for the duration of a note.
func getNoteTicker(secondsPerNote float64) *time.Ticker {
	noteDurationString := fmt.Sprintf("%.5f", secondsPerNote) + "s"
	noteDuration, _ := time.ParseDuration(noteDurationString)

	return time.NewTicker(noteDuration)
}

// playNotes is a ticker loop that plays a set of notes in a 2d string array. The done boolean variable is used
// to indicate that the ticker loop should stop. This avoids an infinite from never ending on a separate go routine.
func playNotes(done chan bool, ticker *time.Ticker, player player.NotePlayer, notes [][]string) {

	i := 0
	numNotes := len(notes)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if len(notes[i]) > 0 {
				//don't block successive notes.
				go player.PlayNotes(notes[i])
			}
			i = (i + 1) % numNotes
		}
	}
}

// handleApplicationError determines how to handle an application level error that was unhandled and bubbled up to the
// top of the stack. Accept an interface because the error may be a panic that we trapped.
func handleApplicationError(err interface{}) {

	var errString string

	stackTrace := string(debug.Stack())

	switch x := err.(type) {
	case string:
		errString = x
	case error:
		errString = x.Error()
	default:
		errString = fmt.Sprintf("%v", x)
	}

	errString += ": Debug Stacktrace: " + stackTrace

	fmt.Printf("Sorry! :( An error occurred in the application %s", errString)
	os.Exit(1)
}
