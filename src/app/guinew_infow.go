//go:build guinew

package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func guinew_infowindow_loop(w *app.Window) error {
	var closeBtn widget.Clickable
	var list widget.List

	list.Axis = layout.Vertical

	theme := material.NewTheme()

	var ops op.Ops

	for e := w.Event(); ; e = w.Event() {
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if closeBtn.Clicked(gtx) {
				w.Perform(system.ActionClose)
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.List(theme, &list).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
						label := material.Body1(theme, item_show_info_get_text())
						label.Alignment = text.Start
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, label.Layout)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(theme, &closeBtn, "Close")
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, btn.Layout)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
