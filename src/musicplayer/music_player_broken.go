//go:build gui

package music_player

type MusicPlayerBroken struct {
	MusicPlayerInterface
}

func (p *MusicPlayerBroken) Play()  {}
func (p *MusicPlayerBroken) Pause() {}
func (p *MusicPlayerBroken) IsPlaying() bool {
	return false
}
