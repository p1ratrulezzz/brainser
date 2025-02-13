//go:build gui

package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"jetbrainser/src/musicplayer"
	"jetbrainser/src/patchers"
	"math/rand"
	"strings"
	"time"
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
	wndMain := a.NewWindow("Jetbrainser" + windowsTitleSuffix)
	wndMain.Resize(fyne.NewSize(640, 480))
	wndMain.SetFixedSize(true)
	wndMain.CenterOnScreen()

	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()

	var appdataDirs []string
	var files []string
	var products []patchers.ProductInfo
	productsMap := make(map[string]*patchers.ProductInfo)
	var guiPrepareExit func()

	files = []string{"test"}
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

	bndWidgetText := binding.NewString()
	wdgListText := widget.NewLabelWithData(bndWidgetText)
	wdgList.Hide()

	wdgLabelTop := widget.NewLabel("Searching for products...")
	wdgProductListCheckbox := widget.NewCheckGroup([]string{}, func(i []string) {})
	wdgProductListCheckbox.Hide()

	step := 0
	selectedIndex := 0
	var ptrWdgButtonNext *widget.Button
	var selectedSource, selectedAppdata, selectedKey int
	wdgButtonNext := widget.NewButton("Patch selected", func() {
		wdgList.UnselectAll()

		switch step {
		// Automatic patch mode
		case 0:
			step = 100
			var productsSelected []patchers.ProductInfo
			for _, key := range wdgProductListCheckbox.Selected {
				productsSelected = append(productsSelected, *productsMap[key])
			}

			messages := doAutoPatch(tool, productsSelected)
			wdgListText.SetText(strings.Join(messages, "\n"))
			guiPrepareExit()
			break
		// Source file selected
		case 1:
			ptrWdgButtonNext.Disable()
			bndData = appdataDirs
			bndFiles.Reload()
			selectedSource = selectedIndex
			wdgLabelTop.SetText("Will patch " + files[selectedSource] + "\nSelect destination directory")
			break
		case 2:
			selectedAppdata = selectedIndex
			wdgLabelTop.SetText("Will use " + appdataDirs[selectedAppdata] + "\nSelect key")
			bndData = KeyListNameIndexed
			bndFiles.Reload()
			ptrWdgButtonNext.SetText("Patch")

			break
		case 3:
			selectedKey = selectedIndex
			wdgLabelTop.SetText("Patched with " + KeyListSlugIndexed[selectedKey])
			guiPrepareExit()
			appdata := appdataDirs[selectedAppdata]
			if selectedAppdata == 0 {
				appdata = ""
			}

			errorMessages := doPatch(files[selectedSource], appdata, "", selectedKey)
			if len(errorMessages) > 0 {
				wdgLabelTop.SetText("Errors occured:" + strings.Join(errorMessages, "\n"))
			} else {
				wdgLabelTop.SetText("Patched!")
			}
			step = 100
			break
		default:
			wndMain.Close()
		}

		step++
	})

	ptrWdgButtonNext = wdgButtonNext

	var wdgButtonRescanPtr *widget.Button
	var wdgButtonManualPtr *widget.Button
	wdgButtonManual := widget.NewButton("Switch to manual mode", func() {
		wdgButtonRescanPtr.Hide()
		wdgList.Show()
		wdgListText.Hide()
		wdgButtonManualPtr.Hide()
		wdgButtonNext.Disable()

		wdgProgressBar.Start()
		wdgProgressBar.Show()

		wdgButtonNext.SetText("Next")
		wdgLabelTop.SetText("Searching for *.vmoptions files")

		step = 1
		filesChan := make(chan []string)
		appdataDirsChan := make(chan []string)
		go guiFindDirectories(tool, filesChan, appdataDirsChan)
		stopTimer := make(chan bool, 1)
		go guiMessagesProgress(wdgLabelTop, stopTimer)

		files = <-filesChan
		appdataDirs = <-appdataDirsChan
		stopTimer <- true

		bndData = files
		bndFiles.Reload()
		wdgProgressBar.Stop()
		wdgProgressBar.Hide()
		wdgLabelTop.SetText(fmt.Sprintf("Found %d files", len(files)))
	})

	wdgButtonManualPtr = wdgButtonManual

	wdgButtonRescan := widget.NewButton("Rescan", func() {
		wdgProductListCheckbox.Hide()
		wdgProductListCheckbox.Options = []string{}
		bndData = []string{}
		bndFiles.Reload()

		wdgProgressBar.Start()
		wdgProgressBar.Show()
		productsChan := make(chan []patchers.ProductInfo)
		go guiFindRunningProducts(tool, productsChan)

		products = <-productsChan
		close(productsChan)

		if len(products) == 0 {
			wdgLabelTop.SetText("No products were found")
			bndWidgetText.Set("No running products were found. \n Make sure you have running products and press \"Rescan\" or use manual mode")
			wdgButtonManual.Enable()
			return
		}

		var productList []string
		for _, productInfo := range products {
			listItem := productInfo.ProductName + "(" + productInfo.VmoptionsSourcePath + ")"
			productList = append(productList, listItem)
			wdgProductListCheckbox.Append(listItem)
			productsMap[listItem] = &productInfo
		}

		wdgProductListCheckbox.SetSelected(productList)

		wdgProductListCheckbox.Show()

		wdgProgressBar.Stop()
		wdgProgressBar.Hide()
		wdgButtonManual.Enable()
		wdgButtonNext.Enable()
		wdgLabelTop.SetText(fmt.Sprintf("Found %d running products", len(productList)))
	})

	guiPrepareExit = func() {
		step = 100
		wdgButtonNext.SetText("Exit")
		wdgButtonNext.Enable()
		wdgButtonManual.Disable()
		wdgButtonRescan.Disable()
	}

	wdgButtonRescanPtr = wdgButtonRescan

	wdgButtonManual.Disable()
	wdgButtonNext.Disable()

	wdgList.OnSelected = func(id widget.ListItemID) {
		wdgButtonNext.Enable()
		selectedIndex = id
	}

	wdgButtonInfo := widget.NewButton("Info", func() {
		_ = wdgProductListCheckbox.Selected
		wdgLabel := widget.NewLabel(item_show_info_get_text())
		var wdgPopupModalPtr *widget.PopUp
		btnClosePopup := widget.NewButton("Close", func() {
			wdgPopupModalPtr.Hide()
		})

		content := container.New(layout.NewVBoxLayout(), wdgLabel, btnClosePopup)
		wdgPopupModal := widget.NewModalPopUp(content, wndMain.Canvas())
		wdgPopupModalPtr = wdgPopupModal

		wdgPopupModal.Show()
	})

	var wdgButtonCleanupModeSwitchPtr *widget.Button
	wdgButtonCleanupModeSwitch := widget.NewButton("Cleanup mode: Off", func() {
		globalvarCleanupMode = !globalvarCleanupMode
		cleanupModeSuffix := ": Off"
		if globalvarCleanupMode {
			cleanupModeSuffix = ": On"
		}

		wdgButtonCleanupModeSwitchPtr.SetText("Cleanup mode" + cleanupModeSuffix)
	})

	wdgButtonCleanupModeSwitchPtr = wdgButtonCleanupModeSwitch

	wdgButtonExit := widget.NewButton("Exit", func() {
		wndMain.Close()
		a.Quit()
	})

	top := container.NewVBox(wdgLabelTop, wdgProductListCheckbox, wdgProgressBar)

	// wdgMusicButton := widget.NewButton("Music", func() {})
	wdgMusicButton := addMusicButton()

	buttons1 := container.NewAdaptiveGrid(2, wdgButtonNext, wdgButtonRescan)
	buttons2 := container.NewAdaptiveGrid(2, wdgButtonCleanupModeSwitch, wdgButtonManual)
	buttons3 := container.NewAdaptiveGrid(3, wdgMusicButton, wdgButtonInfo, wdgButtonExit)

	buttonsBox := container.NewVBox(buttons1, buttons2, buttons3)
	list := container.New(layout.NewMaxLayout(), wdgList, wdgListText)
	content := container.New(
		layout.NewBorderLayout(top, buttonsBox, nil, nil),
		top,
		list,
		buttonsBox,
	)

	go func() {
		// time.Sleep(1 * time.Second)
		wdgButtonRescan.OnTapped()
	}()
	wndMain.SetContent(content)
	wndMain.ShowAndRun()
}

func addMusicButton() *widget.Button {
	musicEnabled := false
	var player musicplayer.MusicPlayerInterface
	wdgButtonMusic := widget.NewButton("Music", func() {
		if player == nil {
			player = musicplayer.NewPlayer()
		}

		musicEnabled = !musicEnabled
		if musicEnabled {
			player.Play()
		} else {
			player.Pause()
		}
	})

	go func() {
		wdgButtonMusic.OnTapped()

		pos := 0
		text := strings.ToLower(wdgButtonMusic.Text)
		for true {
			time.Sleep(300 * time.Millisecond)
			if !musicEnabled {
				continue
			}

			if pos >= len(text) {
				pos = 0
			}

			newText := text[0:pos] + strings.ToUpper(text[pos:pos+1]) + text[pos+1:]
			wdgButtonMusic.SetText(newText)
			pos++
		}
	}()

	return wdgButtonMusic
}

func guiFindDirectories(tool patchers.PatcherTool, files chan []string, appdataDirs chan []string) {
	files <- tool.FindVmoptionsFiles()
	close(files)

	appdataDirsTmp := tool.FindConfigDirectories()
	appdataDirsTmp = append([]string{"Patch in place"}, appdataDirsTmp...)

	appdataDirs <- appdataDirsTmp
	close(appdataDirs)
}

func guiFindRunningProducts(tool patchers.PatcherTool, products chan []patchers.ProductInfo) {
	products <- tool.FindVmoptionsFromProcesses()
}

var messages []string = []string{
	"Waiting...",
	"Still waiting...",
	"This is taking soooo long....",
	"Probably you should use an SSD?",
	"Wow your computer is a potato...",
	"Seriously... still running???",
}

func guiMessagesProgress(label *widget.Label, stopTimer chan bool) {
	currentId := 0
	waitDuration, _ := time.ParseDuration("5s")
	stopTimer <- false
	for !<-stopTimer {
		if currentId < len(messages) {
			label.SetText(messages[currentId])
			currentId++
		}

		time.Sleep(waitDuration)
		duration, _ := time.ParseDuration(fmt.Sprintf("%ds", rand.Intn(9)+1))
		waitDuration += duration
		stopTimer <- false
	}

	label.SetText("")
	close(stopTimer)
}
