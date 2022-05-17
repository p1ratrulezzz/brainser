//go:build gui

package main

import (
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"jetbrainser/src/patchers"
	"strings"
)

func main() {
	flgNogui := flag.Bool("nogui", false, "Disable gui")
	flag.Parse()

	if *flgNogui == true {
		menu_loop()
	} else {
		gui()
	}
}

func gui() {
	a := app.New()
	wndMain := a.NewWindow("Jetbrainser")
	wndMain.Resize(fyne.NewSize(640, 480))
	wndMain.SetFixedSize(true)
	wndMain.CenterOnScreen()

	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()

	var appdataDirs []string
	var files []string
	bndData := files
	bndFiles := binding.BindStringList(&bndData)

	wdgProgressBar := widget.NewProgressBarInfinite()

	wdgList := widget.NewListWithData(bndFiles,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			wdbLbl := o.(*widget.Label)
			wdbLbl.Alignment = fyne.TextAlignTrailing
			wdbLbl.Bind(i.(binding.String))

		})

	wdgLabelTop := widget.NewLabel("Select source *.vmoptions file from the list")

	step := 0
	selectedIndex := 0
	var ptrWdgButtonNext *widget.Button
	var selectedSource, selectedAppdata, selectedKey int
	wdgButtonNext := widget.NewButton("Next", func() {
		wdgList.UnselectAll()

		switch step {
		// Source file selected
		case 0:
			ptrWdgButtonNext.Disable()
			bndData = appdataDirs
			bndFiles.Reload()
			selectedSource = selectedIndex
			wdgLabelTop.SetText("Will patch " + files[selectedSource] + "\nSelect destination directory")
			break
		case 1:
			selectedAppdata = selectedIndex
			wdgLabelTop.SetText("Will use " + appdataDirs[selectedAppdata] + "\nSelect key")
			bndData = KeyList
			bndFiles.Reload()
			ptrWdgButtonNext.SetText("Patch")

			break
		case 2:
			selectedKey = selectedIndex
			wdgLabelTop.SetText("Patched with " + KeyList[selectedKey])
			ptrWdgButtonNext.SetText("Exit")
			appdata := appdataDirs[selectedAppdata]
			if selectedAppdata == 0 {
				appdata = ""
			}

			errorMessages := doPatch(files[selectedSource], appdata, selectedKey)
			if len(errorMessages) > 0 {
				wdgLabelTop.SetText("Errors occured:" + strings.Join(errorMessages, "\n"))
			} else {
				wdgLabelTop.SetText("Patched!")
			}
			break
		default:
			wndMain.Close()
		}

		step++
	})

	ptrWdgButtonNext = wdgButtonNext

	wdgButtonNext.Disable()

	wdgList.OnSelected = func(id widget.ListItemID) {
		wdgButtonNext.Enable()
		selectedIndex = id
	}

	go func() {
		files = tool.FindVmoptionsFiles()

		appdataDirs = tool.FindConfigDirectories()
		appdataDirs = append([]string{"Patch in place"}, appdataDirs...)

		bndData = files
		bndFiles.Reload()

		wdgProgressBar.Stop()
		wdgProgressBar.Hide()

		if len(files) == 0 {
			wdgLabelTop.SetText("No *.vmoptions files were found")
			step = -1
			wdgButtonNext.Enable()
			wdgButtonNext.SetText("Exit")
		}
	}()

	top := container.NewVBox(wdgLabelTop, wdgProgressBar)
	content := container.New(
		layout.NewBorderLayout(top, wdgButtonNext, nil, nil),
		top,
		wdgButtonNext,
		container.New(layout.NewMaxLayout(), wdgList),
	)

	wndMain.SetContent(content)
	wndMain.ShowAndRun()
}
