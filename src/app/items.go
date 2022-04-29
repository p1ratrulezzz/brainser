package main

import (
	"fmt"
	"jetbrainser/src/patchers"
	"path/filepath"
)

func item_show_info() {
	info := `
This tool is provided specially for Typical Society forum and shipped AS IS.
This tool is open source and can be used by anyone. 
Please note, this tool is only for demonstrating and non commercial use, 
if you love some software - buy it and use freely.

Visit: https://forum.losper.net/

`
	print(info)
}

func item_patch() {
	_, keylist := getKeys()

	patcher := patchers.Patcher{osName, nil}
	fmt.Println("Searching for *.vmoptions files ... ")
	files, appdataDirs := patcher.GetTool().FindDirectories()

	if len(files) == 0 {
		fmt.Println("No *.vmoptions files found")
		return
	}

	fmt.Println("Choose what to patch:")
	selected := inputselect_from_array(files)

	fmt.Printf("I will patch %s\n", files[selected])

	appdataDirs = append([]string{"Patch in place"}, appdataDirs...)
	fmt.Println("Select where to put key and *.vmoptions file")
	selectedFolder := inputselect_from_array(appdataDirs)

	appdataSelected := ""
	if selectedFolder == 0 {
		fmt.Printf("I will put key into %s\n", filepath.Dir(files[selected]))
	} else {
		appdataSelected = appdataDirs[selectedFolder]
		fmt.Printf("I will put key into %s\n", appdataDirs[selectedFolder])
	}

	fmt.Println("Choose the key to use")
	chosenKey := inputselect_from_array(keylist)

	doPatch(files[selected], appdataSelected, chosenKey)
}
