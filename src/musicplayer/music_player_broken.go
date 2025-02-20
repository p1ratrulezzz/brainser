//go:build gui

package musicplayer

type MusicPlayerBroken struct {
	MusicPlayerInterface
}

func (p *MusicPlayerBroken) SetFileBytes(fileBytes []byte) {}
func (p *MusicPlayerBroken) Play()                         {}
func (p *MusicPlayerBroken) Pause()                        {}
func (p *MusicPlayerBroken) IsPlaying() bool {
	return false
}
