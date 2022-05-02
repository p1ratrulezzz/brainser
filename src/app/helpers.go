package main

import (
	"fmt"
	"jetbrainser/src/cryptor"
	"log"
	"strconv"
	"strings"
)

var KeyList = []string{
	"appcode",
	"clion",
	"datagrip",
	"goland",
	"idea",
	"phpstorm",
	"pycharm",
	"rider",
	"rubymine",
	"webstorm",
}

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

func getResource(path string) []byte {
	path = strings.TrimSuffix(path, ".enc")
	encrypted, _ := resources.ReadFile("resources_enc/" + path + ".enc")
	resData := cryptor.Decrypt(encrypted)

	return resData
}

func cleanupVmoptions(vmoptionsContent []byte) string {
	vmoptionsContentString := string(vmoptionsContent)
	offset := 0
	needle := "-javaagent:"
	for pos := 0; offset < len(vmoptionsContentString) && pos != -1; {
		pos = strings.Index(vmoptionsContentString[offset:], needle)
		if pos == -1 {
			continue
		}

		offsettmp := offset
		offset += pos + len(needle)
		pos += offsettmp

		if vmoptionsContentString[(pos-1):pos] != "#" {
			vmoptionsContentString = vmoptionsContentString[0:pos] + "#" + vmoptionsContentString[pos:]
			offset--
		}
	}

	if vmoptionsContentString[len(vmoptionsContentString)-1:] != "\n" {
		vmoptionsContentString += "\n"
	}

	return vmoptionsContentString
}
