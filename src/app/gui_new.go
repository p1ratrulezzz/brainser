//go:build guinew

package main

import (
	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"flag"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/colornames"
	"image/color"
	"jetbrainser/src/musicplayer"
	"jetbrainser/src/patchers"
	"log"
	"os"
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
	Product  patchers.ProductInfo
	Checkbox *widget.Bool
}

type AppRes struct {
	IsRescanInProgress            bool
	btnPatch, btnRescan, btnMusic widget.Clickable
	Button1, Button2, Button3     widget.Clickable
	ProductInfoCheckboxItems      []ProductInfoCheckbox
	Tool                          patchers.PatcherTool
	Player                        musicplayer.MusicPlayerInterface
	MusicIsPlaying                bool
	W                             *app.Window
}
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
		W:      w,
	}

	res.Player.SetFileBytes(getGorchichka())

	go func() {
		handlerWindowOnLoad(&res)
	}()

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
								guinewCheckboxesChildren(&res)...,
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
										btn := material.Button(th, &res.btnRescan, "Rescan")
										if res.IsRescanInProgress {
											btn.Background = color.NRGBA(colornames.Gray)
										}

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

func guinewCheckboxesChildren(res *AppRes) []layout.FlexChild {
	var children []layout.FlexChild
	for _, item := range res.ProductInfoCheckboxItems {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.CheckBox(material.NewTheme(), item.Checkbox, item.Product.ProductSlug).Layout(gtx)
		}))
	}

	return children
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
		go handlerBtnRescanClick(res)
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

func handlerBtnRescanClick(res *AppRes) {
	if res.IsRescanInProgress {
		return
	}

	res.IsRescanInProgress = true
	res.ProductInfoCheckboxItems = []ProductInfoCheckbox{}

	productsChan := make(chan []patchers.ProductInfo)
	go func() {
		productsChan <- res.Tool.FindVmoptionsFromProcesses()
	}()
	products := <-productsChan
	close(productsChan)

	for _, product := range products {
		item := ProductInfoCheckbox{
			Product:  product,
			Checkbox: new(widget.Bool),
		}

		res.ProductInfoCheckboxItems = append(res.ProductInfoCheckboxItems, item)
	}

	res.IsRescanInProgress = false
	res.W.Invalidate()
}
