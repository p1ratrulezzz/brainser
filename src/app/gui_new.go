//go:build guinew

package main

import (
	"flag"
	"fmt"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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
			app.Title("Centered Window with Modal"),
			app.Size(unit.Dp(800), unit.Dp(600)), // Устанавливаем размер окна
		)

		err := loop(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type AppRes struct {
	Button1, Button2, Button3 widget.Clickable
	CloseModal                widget.Clickable
	Items                     []string
	ShowModal                 bool
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	res := AppRes{
		Items: []string{"Элемент 1", "Элемент 2", "Элемент 3", "Элемент 4", "Элемент 5"},
	}
	// Список (пример данных)
	list := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}

	for e := w.Event(); ; e = w.Event() {
		switch e := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Обработка нажатий кнопок
			if res.Button1.Clicked(gtx) {
				res.ShowModal = true // Показать модальное окно
			}
			if res.Button2.Clicked(gtx) {
				fmt.Println("Нажата кнопка 2")
			}
			if res.Button3.Clicked(gtx) {
				fmt.Println("Нажата кнопка 3")
			}
			if res.CloseModal.Clicked(gtx) {
				res.ShowModal = false // Закрыть модальное окно
			}

			// Главная компоновка с наложением модального окна
			layout.Stack{}.Layout(gtx,
				// Основной слой (фон)
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceEnd,
					}.Layout(gtx,
						// Верхняя часть: 80% для списка
						layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
							return material.List(th, &list).Layout(gtx, len(res.Items), func(gtx layout.Context, index int) layout.Dimensions {
								return material.Body1(th, res.Items[index]).Layout(gtx)
							})
						}),
						// Нижняя часть: 20% для кнопок
						layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceEvenly,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &res.Button1, "Показать модалку")
									return btn.Layout(gtx)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &res.Button2, "Кнопка 2")
									return btn.Layout(gtx)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &res.Button3, "Кнопка 3")
									return btn.Layout(gtx)
								}),
							)
						}),
					)
				}),
				// Модальное окно как наложение
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					if !res.ShowModal {
						return layout.Dimensions{} // Не показывать, если модалка закрыта
					}

					// Затемнение фона
					paint.FillShape(gtx.Ops, color.NRGBA{A: 128}, clip.Rect{Max: gtx.Constraints.Max}.Op())

					// Модальное окно в центре
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(300) // Ширина модального окна
						gtx.Constraints.Max.Y = gtx.Dp(200) // Высота модального окна
						return widget.Border{
							Color: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
							Width: unit.Dp(1),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(material.H6(th, "Модальное окно").Layout),
									layout.Rigid(material.Body1(th, "Это пример модального окна.").Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th, &res.CloseModal, "Закрыть")
										return btn.Layout(gtx)
									}),
								)
							})
						})
					})
				}),
			)

			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}
