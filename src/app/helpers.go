package main

import (
	"fmt"
	"io/fs"
	"log"
	"strconv"
)

func delay() {
	fmt.Print("Press 'Enter' to continue...")
	stdin.ReadBytes('\n')
}

func inputselect_from_array(choses []string) int {
	if len(choses) == 0 {
		log.Fatal("empty list passed")
	}

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
	keys, _ := resources.ReadDir("resources/keys")

	for _, entry := range keys {
		list = append(list, entry.Name())
	}

	return keys, list
}
