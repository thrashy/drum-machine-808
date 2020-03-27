package player

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// WavNotePlayer defines a struct for playing notes using a wav player.
type WavNotePlayer struct {
	FileMap map[string]string
}

// NewWavNotePlayer initializes a WavNotePlayer.
func NewWavNotePlayer() *WavNotePlayer {

	format := beep.Format{
		NumChannels: 1,
		SampleRate:  44100,
		Precision:   2,
	}

	// Calling speaker.init multiple times will reset the speaker, preventing multiple sounds from playing
	// simultaneously. See: https://github.com/faiface/beep/wiki/Hello,-Beep!
	// Ignore the error returned. It is not used in any of the examples.
	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// See: http://smd-records.com/tr808/?page_id=14 for details of the sounds, volumes and decay values.
	return &WavNotePlayer{
		FileMap: map[string]string{
			"bd": "BD/BD7550",
			"sn": "SD/SD5050",
			"hh": "CH/CH",
			"oh": "OH/OH25",
			"rd": "CY/CY2550",
		},
	}
}

// PlayNotes uses standard output to print the names of notes being played.
func (player *WavNotePlayer) PlayNotes(noteNames []string) interface{} {

	const folderPath = "../assets/sounds/tr808wav/"
	const extension = ".WAV"

	streamClosers := make([]*beep.StreamSeekCloser, 0)

	defer player.closeStreamers(streamClosers)

	for _, note := range noteNames {

		file, exists := player.FileMap[note]

		if exists {
			filePath := folderPath + file + extension
			// TODO: handle error
			streamCloser, _ := player.getStreamer(filePath)
			streamClosers = append(streamClosers, streamCloser)
		}
	}

	if len(streamClosers) > 0 {
		player.playStreams(streamClosers)
	}

	return nil
}

// playNote plays a single note with a relative file path used to find the sound wav file.
func (player *WavNotePlayer) playStreams(streamClosers []*beep.StreamSeekCloser) {

	// There should be a better way to do this. Need to do a type conversion from []beep.StreamSeekCloser
	// to []beep.Streamer
	streamers := make([]beep.Streamer, 0)
	for _, streamer := range streamClosers {
		streamers = append(streamers, *streamer)
	}
	// beep handles executing these concurrently
	speaker.Play(beep.Mix(streamers...))

	// The beep callback doesn't seem to be working. Add an artificial time to block. Need more time too look into
	// this. Add an artificial time wait. This is all executed on a separate go routine, so we won't block the next
	// note.
	time.Sleep(2 * time.Second)
}

// getStreamer gets a streamer for a particular wav file.
func (player *WavNotePlayer) getStreamer(filePath string) (*beep.StreamSeekCloser, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	streamer, _, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}

	return &streamer, nil
}

// closeStreamers closes all the open streams.
func (player *WavNotePlayer) closeStreamers(streams []*beep.StreamSeekCloser) {
	for _, streamer := range streams {
		// return from Close() is ignored in examples.
		_ = (*streamer).Close()
	}
}
