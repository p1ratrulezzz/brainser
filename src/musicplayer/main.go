//go:build gui || guinew

package musicplayer

import (
	"bytes"
	"errors"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"io"
	"time"
)

type MusicPlayerInterface interface {
	Play()
	Pause()
	IsPlaying() bool
	SetFileBytes(fileBytes []byte)
}

type MusicPlayer struct {
	MusicPlayerInterface
	otoCtx        *oto.Context
	currentPlayer *oto.Player
	currentFile   *mp3.Decoder
	fileBytes     []byte
}

func (player *MusicPlayer) SetFileBytes(fileBytes []byte) {
	player.fileBytes = fileBytes
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
	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(player.fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		return
	}

	player.currentFile = decodedMp3
	player.currentPlayer = player.otoCtx.NewPlayer(player.currentFile)
	player.currentPlayer.SetVolume(0.35)
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
	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)

	if err != nil {
		player := new(MusicPlayerBroken)
		return player
	}

	<-readyChan

	player := new(MusicPlayer)
	player.otoCtx = otoCtx
	player.currentPlayer = nil

	return player
}
