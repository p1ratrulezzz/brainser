package main

import (
	"fmt"
	"jetbrainser/src/patchers"
	"path/filepath"
)

func item_patch() {
	fmt.Println("Before patching, start your JetBrains product once and close it. Then press enter to continue...")
	stdin.ReadLine()

	patcher := patchers.Patcher{osName, nil}
	fmt.Println("Searching for *.vmoptions files ... ")
	tool := patcher.GetTool()
	files := tool.FindVmoptionsFiles()

	var sourceVmoptionsPath string
	var appdataSelected string
	if len(files) > 0 {
		fmt.Println("Choose what to patch:")
		selected := inputselect_from_array(files)
		sourceVmoptionsPath = files[selected]
		fmt.Printf("I will patch %s\nSearching for config dirs...", sourceVmoptionsPath)

		appdataDirs := patcher.GetTool().FindConfigDirectories()
		appdataDirs = append([]string{"Patch in place"}, appdataDirs...)
		fmt.Println("Select where to put key and *.vmoptions file")
		fmt.Println("If it is a Toolbox installed software, choose 0 - Patch in place")
		selectedFolder := inputselect_from_array(appdataDirs)
		if selectedFolder != 0 {
			appdataSelected = appdataDirs[selectedFolder]
			fmt.Printf("I will put key into %s\n", appdataSelected)
		}

	} else {
		fmt.Println("No *.vmoptions files found in program files dir")
		fmt.Println("Searching for *.vmoptions in config dir")

		filesInConfigDir := patcher.GetTool().FindVmoptionsFilesInConfigDir()
		if len(filesInConfigDir) == 0 {
			fmt.Println("No *.vmoptions files found in config dir. Program can't be continued")
			return
		}

		fmt.Println("Choose a file to patch:")
		selected := inputselect_from_array(filesInConfigDir)
		sourceVmoptionsPath = filesInConfigDir[selected]
		appdataSelected = ""
	}

	if appdataSelected == "" {
		fmt.Printf("I will put key into %s\n", filepath.Dir(sourceVmoptionsPath))
	}

	fmt.Println("Choose the key to use (by your product type)")
	chosenKeyIndex := inputselect_from_array(KeyList)

	doPatch(sourceVmoptionsPath, appdataSelected, chosenKeyIndex)
}
