package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// SongBeatConfiguration is a model container for song beat definitions.
func TestUnmarshalTOML(t *testing.T) {

	validTOML := `
		[[song_beat]]
			name = "Four on the Floor"
			beats_per_sequence = 4
			notes_per_beat = 1
			beats_per_minute = 128

			[song_beat.pattern]
				bd = [ 1, 0, 1, 0 ]
				sn = [ 0, 0, 1, 0 ]
				hh = [ 0, 1, 0, 1 ]
		`

	sb := SongBeatConfiguration{}
	err := sb.UnmarshalTOML([]byte(validTOML))

	if err != nil {
		t.Errorf("Error returned from UnmarshalTOML %s.", err.Error())
	}

	if sb.SongBeatData == nil {
		t.Errorf("SongBeatData empty.")
	}

	songBeatData := sb.SongBeatData["Four on the Floor"]

	if songBeatData == nil {
		t.Errorf("SongBeat empty.")
	} else {

		if songBeatData.BeatsPerSequence != 4 {
			t.Errorf("BeatsPerSequence did not match, actual: %d, expected: %d.", songBeatData.BeatsPerSequence, 4)
		}

		if songBeatData.NotesPerBeat != 1 {
			t.Errorf("NotesPerBeat did not match, actual: %d, expected: %d.", songBeatData.NotesPerBeat, 1)
		}

		if songBeatData.BeatsPerMinute != 128 {
			t.Errorf("BeatsPerMinute did not match, actual: %d, expected: %d.", songBeatData.BeatsPerMinute, 128)
		}

		if songBeatData.BeatPattern == nil {
			t.Errorf("BeatPattern empty.")
		}

		if len(songBeatData.BeatPattern) != 3 {
			t.Errorf("BeatPattern length did not match, actual: %d, expected: %d.", len(songBeatData.BeatPattern), 3)
		}

		if songBeatData.BeatPattern["sn"] == nil {
			t.Errorf("BeatPattern for 'snare drum'(sn) empty")
		}
	}

	// cmp is a safer alternative to reflect.DeepEqual for comparing whether two values are semantically equal.
	// See: https://github.com/google/go-cmp for more information.
	if songBeatData != nil && !cmp.Equal(songBeatData.BeatPattern["sn"], []bool{false, false, true, false}) {
		t.Errorf("BeatPattern for 'snare drum'(sn) did not match, actual: %v, expected: %v.",
			songBeatData.BeatPattern["sn"], []bool{false, false, true, false})
	}
}

// SongBeatConfigurationValidateTest tests the Validate function of SongBeatConfiguration
func TestSongBeatConfigurationValidate(t *testing.T) {

	sb := setupSongBeatConfigurationData()

	actual := sb.Validate()

	if actual != nil {
		t.Errorf("Error should have been nil, actual: %s", actual)
	}

	sb.SongBeatData["test"].BeatPattern["bd"] = []bool{true}
	actual = sb.Validate()

	if actual == nil {
		t.Errorf("Error expected but found nil")
	}

	sb = setupSongBeatConfigurationData()
	sb.SongBeatData["test"].BeatsPerMinute = 0
	actual = sb.Validate()

	if actual == nil {
		t.Errorf("Error expected but found nil")
	}

	sb = setupSongBeatConfigurationData()
	sb.SongBeatData["test"].BeatsPerSequence = 0
	actual = sb.Validate()

	if actual == nil {
		t.Errorf("Error expected but found nil")
	}

	sb = setupSongBeatConfigurationData()
	sb.SongBeatData["test"].NotesPerBeat = 0
	actual = sb.Validate()

	if actual == nil {
		t.Errorf("Error expected but found nil")
	}

	sb = setupSongBeatConfigurationData()
	sb.SongBeatData[""] = &SongBeat{
		BeatsPerSequence: 4,
		NotesPerBeat:     2,
		BeatsPerMinute:   128,
		BeatPattern: map[string][]bool{
			"bd": {
				true,
				false,
				true,
				true,
			},
		},
	}
	actual = sb.Validate()

	if actual == nil {
		t.Errorf("Error expected but found nil")
	}
}

// setupSongBeatConfigurationData initializes a valid SongBeatConfiguration.
func setupSongBeatConfigurationData() SongBeatConfiguration {
	return SongBeatConfiguration{
		SongBeatData: map[string]*SongBeat{
			"test": {
				BeatsPerSequence: 2,
				NotesPerBeat:     2,
				BeatsPerMinute:   128,
				BeatPattern: map[string][]bool{
					"bd": {
						true,
						false,
						true,
						true,
					},
				},
			},
		},
	}
}
