//go:build guinew

package main

import (
	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"flag"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"jetbrainser/src/musicplayer"
	"jetbrainser/src/patchers"
)

func main() {
	checkIntegrity()

	flgNogui := flag.Bool("nogui", false, "Disable gui")
	flag.Parse()

	if *flgNogui == true {
		menu_loop()
	} else {
		guinew()
	}
}

func guinew() {
	go func() {
		window := new(app.Window)
		window.Option(
			app.Title("Jetbrainser"+windowsTitleSuffix),
			app.Size(unit.Dp(640), unit.Dp(480)),
		)

		err := loop(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type ProductInfoCheckbox struct {
	Product patchers.ProductInfo
}

type AppRes struct {
	btnPatch, btnRescan, btnMusic widget.Clickable
	Button1, Button2, Button3     widget.Clickable
	ProductInfoCheckboxItems      []ProductInfoCheckbox
	Tool                          patchers.PatcherTool
	Player                        musicplayer.MusicPlayerInterface
	MusicIsPlaying                bool
}

type GioUiCb func(C) D

type (
	D = layout.Dimensions
	C = layout.Context
)

func loop(w *app.Window) error {
	th := material.NewTheme()
	// Load a color emoji font.
	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err == nil {
		collection := gofont.Collection()
		th.Shaper = text.NewShaper(text.WithCollection(append(collection, faces...)))
	}

	var ops op.Ops

	res := AppRes{
		Tool:   createPatcherTool(),
		Player: musicplayer.NewPlayer(),
	}

	res.Player.SetFileBytes(getGorchichka())

	go func() {
		handlerWindowOnLoad(&res)
	}()

	// checkboxGroup := new(widget.Enum)
	var checkbox widget.Bool
	var checkbox2 widget.Bool

	for e := w.Event(); ; e = w.Event() {
		switch e := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			guinewBtnRescanRedraw(&res, gtx)
			guinewBtnMusicRedraw(&res, gtx)

			layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceEnd,
					}.Layout(gtx,
						layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceEnd,
							}.Layout(gtx,
								// Первый чекбокс
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									cb := material.CheckBox(th, &checkbox, "Опция 1")
									return cb.Layout(gtx)
								}),
								// Второй чекбокс
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									cb := material.CheckBox(th, &checkbox2, "Опция 2")
									return cb.Layout(gtx)
								}),
								// Третий чекбокс
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									cb := material.CheckBox(th, &checkbox, "Опция 3")
									return cb.Layout(gtx)
								}),
							)
						}),
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceEvenly,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return marginsButton(gtx, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &res.btnPatch, "Patch")
										return btn.Layout(gtx)
									})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return marginsButton(gtx, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &res.btnMusic, "Music")
										btn.Background = color.NRGBA(colornames.Gray)
										if res.MusicIsPlaying {
											btn.Background = color.NRGBA(colornames.Green)
											btn.Text = "🔊" + btn.Text
										} else {
											btn.Text = "🔇" + btn.Text
										}

										return btn.Layout(gtx)
									})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return marginsButton(gtx, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &res.Button3, "Кнопка 3")
										return btn.Layout(gtx)
									})
								}),
							)
						}),
					)
				}),
			)

			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func marginsButton(gtx C, wg layout.Widget) D {
	return layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		Right:  unit.Dp(35),
		Left:   unit.Dp(35),
	}.Layout(gtx, wg)
}
func guinewBtnRescanRedraw(res *AppRes, gtx C) {
	if res.btnRescan.Clicked(gtx) {

	}
}

func guinewBtnMusicRedraw(res *AppRes, gtx C) {
	if res.btnMusic.Clicked(gtx) {
		handlerBtnMusicClick(res)
	}
}

func createPatcherTool() patchers.PatcherTool {
	patcher := patchers.Patcher{osName, nil, getPomidori()}
	return patcher.GetTool()
}

func rescanProcesses() {}

func handlerWindowOnLoad(res *AppRes) {
	handlerBtnMusicClick(res)
}

func handlerBtnMusicClick(res *AppRes) {
	if !res.MusicIsPlaying {
		res.Player.Play()
	} else {
		res.Player.Pause()
	}

	res.MusicIsPlaying = !res.MusicIsPlaying
}
