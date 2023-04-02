//go:build gui

package music_player

import (
	"bytes"
	"embed"
	"errors"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"io"
	"time"
)

//go:embed music
var musicFiles embed.FS

type MusicPlayerInterface interface {
	Play()
	Pause()
	IsPlaying() bool
}

type MusicPlayer struct {
	MusicPlayerInterface
	otoCtx        *oto.Context
	currentPlayer oto.Player
	currentFile   *mp3.Decoder
}

func (player *MusicPlayer) MusicIsEmptyOrFinished() bool {
	if player.currentPlayer == nil {
		return true
	}

	isEof, err := player.IsEof()
	if err != nil {
		return true
	}

	return isEof
}

func (player *MusicPlayer) ReloadFile() {
	fileBytes, err := musicFiles.ReadFile("music/outrun.mp3")
	if err != nil {
		return
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		return
	}

	player.currentFile = decodedMp3
	player.currentPlayer = player.otoCtx.NewPlayer(player.currentFile)
	player.currentPlayer.SetVolume(1)
}

func (player *MusicPlayer) Play() {
	if player.MusicIsEmptyOrFinished() {
		player.ReloadFile()
	}

	if !player.MusicIsEmptyOrFinished() {
		player.currentPlayer.Play()
		player.PlayRoutine()
	}
}

func (player *MusicPlayer) IsEof() (bool, error) {
	if player.currentFile == nil {
		return true, errors.New("no file loaded")
	}

	pos, err := player.currentFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return true, err
	}

	return pos >= player.currentFile.Length(), nil
}

func (player *MusicPlayer) PlayRoutine() {
	go func() {
		for isEof, err := player.IsEof(); !isEof && err == nil; isEof, err = player.IsEof() {
			time.Sleep(time.Second)
		}

		player.Play()
	}()
}

func (player *MusicPlayer) Pause() {
	if player.MusicIsEmptyOrFinished() {
		return
	}

	player.currentPlayer.Pause()
}

func (player *MusicPlayer) IsPlaying() bool {
	if player.MusicIsEmptyOrFinished() {
		return false
	}

	return player.currentPlayer.IsPlaying()
}

func NewPlayer() MusicPlayerInterface {
	samplingRate := 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	<-readyChan

	if err != nil {
		player := new(MusicPlayerBroken)
		return player
	}

	player := new(MusicPlayer)
	player.otoCtx = otoCtx
	player.currentPlayer = nil

	return player
}

//
//func initMusic() (oto.Player, *mp3.Decoder, error) {
//	fileBytes, err := musicFiles.ReadFile("music/outrun.mp3")
//	if err != nil {
//		return nil, nil, err
//	}
//
//	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
//	fileBytesReader := bytes.NewReader(fileBytes)
//
//	// Decode file
//	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	// Prepare an Oto context (this will use your default audio device) that will
//	// play all our sounds. Its configuration can't be changed later.
//
//	// Usually 44100 or 48000. Other values might cause distortions in Oto
//	samplingRate := 44100
//
//	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
//	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
//	numOfChannels := 2
//
//	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
//	audioBitDepth := 2
//
//	// Remember that you should **not** create more than one context
//	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
//	if err != nil {
//		return nil, nil, err
//	}
//	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
//	<-readyChan
//
//	// Create a new 'player' that will handle our sound. Paused by default.
//	player := otoCtx.NewPlayer(decodedMp3)
//
//	// Play starts playing the sound and returns without waiting for it (Play() is async).
//	return player, decodedMp3, nil
//}
