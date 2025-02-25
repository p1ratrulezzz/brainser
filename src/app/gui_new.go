//go:build guinew

package main

import (
	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"flag"
	"fmt"
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
	"strings"
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
			app.Title("Jetbrainser "+Version+" ("+getKapusta()+") Build: "+BuildNumber+""),
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
	IsInfoWindowOpen              bool
	IsRescanInProgress            bool
	btnPatch, btnRescan, btnMusic widget.Clickable
	btnInfo                       widget.Clickable
	ProductInfoCheckboxItems      []ProductInfoCheckbox
	Tool                          patchers.PatcherTool
	Player                        musicplayer.MusicPlayerInterface
	MusicIsPlaying                bool
	W                             *app.Window
	IsPatchButtonDisabled         bool
	TextLabelPrefix               string
	TextLabelSuffix               string
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
		Tool:            createPatcherTool(),
		Player:          musicplayer.NewPlayer(),
		W:               w,
		TextLabelPrefix: "Press rescan to find running products",
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
			guinewBtnPatchRedraw(&res, gtx)
			guiCheckboxRedraw(&res, gtx)
			guiBtnInfoRedraw(&res, gtx)

			layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceEnd,
					}.Layout(gtx,
						layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {

							var elementsPre []layout.FlexChild
							var elementsAfter []layout.FlexChild

							elementsPre = append(elementsPre, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(20), res.TextLabelPrefix)
								return label.Layout(gtx)
							}))

							checkboxesElements := guinewCheckboxesChildren(&res)
							elementsAfter = append(elementsAfter, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(20), res.TextLabelSuffix)
								return label.Layout(gtx)
							}))

							elements := make([]layout.FlexChild, len(elementsPre)+len(checkboxesElements)+len(elementsAfter))
							pos := copy(elements[0:], elementsPre)
							pos += copy(elements[pos:], checkboxesElements)
							pos += copy(elements[pos:], elementsAfter)

							return layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceEnd,
							}.Layout(gtx,
								elements...,
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
										if res.IsPatchButtonDisabled {
											btn.Background = color.NRGBA(colornames.Gray)
										}

										return btn.Layout(gtx)
									})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return marginsButton(gtx, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &res.btnInfo, "Info")
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
			return material.CheckBox(material.NewTheme(), item.Checkbox, item.Product.ProductName+"("+item.Product.VmoptionsDestinationPath+")").Layout(gtx)
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
		go handlerBtnMusicClick(res)
	}
}

func guinewBtnPatchRedraw(res *AppRes, gtx C) {
	if res.btnPatch.Clicked(gtx) {
		go handlerBtnPatchClick(res)
	}
}

func guiCheckboxRedraw(res *AppRes, gtx C) {
	for _, item := range res.ProductInfoCheckboxItems {
		if item.Checkbox.Update(gtx) {
			go handlerCheckboxClick(res)
			return
		}
	}
}

func guiBtnInfoRedraw(res *AppRes, gtx C) {
	if res.btnInfo.Clicked(gtx) {
		go handlerBtnInfoClick(res)
	}
}

func handlerBtnInfoClick(res *AppRes) {
	if res.IsInfoWindowOpen {
		return
	}

	res.IsInfoWindowOpen = true
	w := new(app.Window)
	w.Option(
		app.Title("Info"),
		app.Size(unit.Dp(300), unit.Dp(300)),
	)

	err := guinew_infowindow_loop(w)
	if err != nil {
		log.Fatal(err)
	}

	res.IsInfoWindowOpen = false
}

func createPatcherTool() patchers.PatcherTool {
	patcher := patchers.Patcher{osName, nil, getPomidori()}
	return patcher.GetTool()
}

func handlerWindowOnLoad(res *AppRes) {
	go handlerBtnMusicClick(res)
	go handlerBtnRescanClick(res)
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

	res.TextLabelSuffix = ""
	res.IsRescanInProgress = true
	res.IsPatchButtonDisabled = true
	res.ProductInfoCheckboxItems = []ProductInfoCheckbox{}
	res.TextLabelPrefix = "Scanning..."

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

		item.Checkbox.Value = true
		res.ProductInfoCheckboxItems = append(res.ProductInfoCheckboxItems, item)
	}

	res.TextLabelPrefix = fmt.Sprintf("Found %d products", len(res.ProductInfoCheckboxItems))
	res.IsRescanInProgress = false
	res.IsPatchButtonDisabled = false

	handlerCheckboxClick(res)

	res.W.Invalidate()
}

func handlerBtnPatchClick(res *AppRes) {
	if res.IsPatchButtonDisabled {
		return
	}

	res.IsPatchButtonDisabled = true
	res.IsRescanInProgress = true

	selectedProducts := gatherSelectedProducts(res)

	if len(selectedProducts) == 0 {
		res.TextLabelSuffix = "No products selected"
		res.IsPatchButtonDisabled = false
		res.IsRescanInProgress = false
		res.W.Invalidate()
		return
	}

	res.TextLabelPrefix = "Patching..."
	res.ProductInfoCheckboxItems = []ProductInfoCheckbox{}

	messages := doAutoPatch(res.Tool, selectedProducts)

	res.TextLabelPrefix = "Done"

	res.TextLabelSuffix = strings.Join(messages, "\n")

	res.IsRescanInProgress = false

	res.W.Invalidate()
}

func gatherSelectedProducts(res *AppRes) []patchers.ProductInfo {
	var selectedProducts []patchers.ProductInfo
	for _, item := range res.ProductInfoCheckboxItems {
		if item.Checkbox.Value {
			selectedProducts = append(selectedProducts, item.Product)
		}
	}

	return selectedProducts
}

func handlerCheckboxClick(res *AppRes) {
	res.IsPatchButtonDisabled = len(gatherSelectedProducts(res)) == 0
}
