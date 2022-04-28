package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
)

func delay() {
	fmt.Print("Press 'Enter' to continue...")
	stdin.ReadBytes('\n')
}

func inputselect_from_array(choses []string) int {
	for {
		for i, label := range choses {
			fmt.Printf("%d: %s\n", i, label)
		}

		inbuf, _, _ := stdin.ReadLine()
		selected, err := strconv.Atoi(string(inbuf))
		if err != nil || selected < 0 || selected >= len(choses) {
			fmt.Println("Incorrect choice. Select correct one")
			continue
		}

		return selected
	}
}

func getKeys() ([]fs.DirEntry, []string) {
	var list []string
	keys, _ := resources.ReadDir(filepath.Join("resources", "keys"))

	for _, entry := range keys {
		list = append(list, entry.Name())
	}

	return keys, list
}
