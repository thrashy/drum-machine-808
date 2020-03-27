package model

import (
	"errors"
	"fmt"

	"github.com/pelletier/go-toml"
)

// SongBeatConfiguration is a model container for song beat definitions.
type SongBeatConfiguration struct {
	SongBeatData map[string]*SongBeat
}

// SongBeat defines a single song within the configuration file.
type SongBeat struct {
	BeatsPerSequence int
	NotesPerBeat     int
	BeatsPerMinute   int
	BeatPattern      map[string][]bool
}

// UnmarshalTOML takes toml configuration data as a byte array and deserializes it into a struct The third party
// package has difficulty with nested structs or non-primitive data types. Ideally this implementation should be a
// single call to toml.Unmarshall.
func (sb *SongBeatConfiguration) UnmarshalTOML(data []byte) error {

	tree, err := toml.LoadBytes(data)

	if err != nil {
		return err
	}

	songBeatData := tree.Get("song_beat").([]*toml.Tree)

	numSongBeats := len(songBeatData)
	songBeats := make(map[string]*SongBeat)

	for i := 0; i < numSongBeats; i++ {

		name := songBeatData[i].Get("name").(string)
		songBeats[name] = &SongBeat{}

		// The toml package uses int64 internally, we need to cast to int
		songBeats[name].BeatsPerSequence = int(songBeatData[i].Get("beats_per_sequence").(int64))
		songBeats[name].NotesPerBeat = int(songBeatData[i].Get("notes_per_beat").(int64))
		songBeats[name].BeatsPerMinute = int(songBeatData[i].Get("beats_per_minute").(int64))

		pattern := songBeatData[i].Get("pattern").(*toml.Tree)

		sb.unmarshalPattern(songBeats, pattern, name)
	}

	sb.SongBeatData = songBeats

	return nil
}

// unmarshalPattern takes the [song_beat.pattern] section from the config toml and deserializes it into a map.
func (sb *SongBeatConfiguration) unmarshalPattern(songBeats map[string]*SongBeat, pattern *toml.Tree, name string) {

	soundKeys := pattern.Keys()
	numSounds := len(soundKeys)
	songBeats[name].BeatPattern = make(map[string][]bool)

	for j := 0; j < numSounds; j++ {
		soundPattern := pattern.Get(soundKeys[j]).([]interface{})
		numNotes := len(soundPattern)
		binData := make([]bool, numNotes, numNotes)

		for k := 0; k < numNotes; k++ {
			//convert binary data to bool, anything non-zero should be treated as true
			binData[k] = int(soundPattern[k].(int64)) != 0
		}

		songBeats[name].BeatPattern[soundKeys[j]] = binData
	}
}

// Validate verifies that the data in a configuration instance is valid. This is useful to avoid possibly more cryptic
// errors further down the stack.
func (sb *SongBeatConfiguration) Validate() error {

	for name, value := range sb.SongBeatData {

		if len(name) == 0 {
			return errors.New("each song must have name")
		}

		if value.BeatsPerMinute == 0 {
			return errors.New("each song must have beats per minute defined")
		}

		if value.BeatsPerSequence == 0 {
			return errors.New("each song must have beats per sequence defined")
		}

		if value.NotesPerBeat == 0 {
			return errors.New("each song must have notes per beat defined")
		}

		numNotes := value.NotesPerBeat * value.BeatsPerSequence

		for beatType, notes := range value.BeatPattern {
			if len(notes) != numNotes {
				return fmt.Errorf("beat pattern '%s' must match NotesPerBeat times BeatsPerSequence and "+
					"have %d notes", beatType, numNotes)
			}
		}
	}

	return nil
}
