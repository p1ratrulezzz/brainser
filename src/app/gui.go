package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"jetbrainser/src/patchers"
)

func gui() {
	a := app.New()
	wndMain := a.NewWindow("Jetbrainser")
	wndMain.Resize(fyne.NewSize(640, 480))
	wndMain.SetFixedSize(true)
	wndMain.CenterOnScreen()

	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()
	files := tool.FindVmoptionsFiles()
	appdataDirs := patcher.GetTool().FindConfigDirectories()
	appdataDirs = append([]string{"Patch in place"}, appdataDirs...)

	bndData := files
	bndFiles := binding.BindStringList(&bndData)

	wdgList := widget.NewListWithData(bndFiles,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	wdgLabelTop := widget.NewLabel("Select source *.vmoptions file from the list")
	step := 0
	selectedIndex := 0
	var ptrWdgButtonNext *widget.Button
	wdgButtonNext := widget.NewButton("Next", func() {
		var selectedSource, selectedAppdata, selectedKey int
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

			doPatch(files[selectedSource], appdata, selectedKey)
			break
		case 3:
			wndMain.Close()
			break
		default:
			panic("unknown step")
		}

		step++
	})

	ptrWdgButtonNext = wdgButtonNext

	wdgButtonNext.Disable()

	wdgList.OnSelected = func(id widget.ListItemID) {
		wdgButtonNext.Enable()
		selectedIndex = id
	}

	content := container.New(
		layout.NewBorderLayout(wdgLabelTop, wdgButtonNext, nil, nil),
		wdgLabelTop,
		container.New(layout.NewMaxLayout(), wdgList),
		wdgButtonNext,
	)

	wndMain.SetContent(content)

	wndMain.ShowAndRun()
}
