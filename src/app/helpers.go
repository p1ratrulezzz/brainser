package main

import (
	"errors"
	"fmt"
	"hash/adler32"
	"jetbrainser/src/cryptor"
	"log"
	"os"
	"strconv"
	"strings"
)

var KeyList = map[string]string{
	"appcode":  "AppCode",
	"clion":    "Clion",
	"datagrip": "Datagrip",
	"goland":   "GoLand",
	"idea":     "Idea",
	"phpstorm": "PhpStorm",
	"pycharm":  "PyCharm",
	"rider":    "Rider",
	"rubymine": "RubyMine",
	"webstorm": "WebStorm",
}

var KeyListSlugIndexed, KeyListNameIndexed []string

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

func cleanupVmoptions(vmoptionsContent []byte) (string, []string) {
	vmoptionsContentString := string(vmoptionsContent)
	offset := 0
	needle := "-javaagent:"
	var agents []string

	for pos := 0; offset < len(vmoptionsContentString) && pos != -1; {
		pos = strings.Index(vmoptionsContentString[offset:], needle)
		if pos == -1 {
			continue
		}

		offsettmp := offset
		offset += pos + len(needle)
		pos += offsettmp

		if vmoptionsContentString[(pos-1):pos] != "#" {
			lineEndPos := strings.Index(vmoptionsContentString[pos:], "\n")
			if lineEndPos == -1 {
				lineEndPos = len(vmoptionsContentString)
			} else {
				lineEndPos += pos
			}

			agent := vmoptionsContentString[pos+len(needle) : lineEndPos]
			agent = strings.Trim(agent, "\n\r\t ")
			agents = append(agents, agent)
			vmoptionsContentString = vmoptionsContentString[0:pos] + "#" + vmoptionsContentString[pos:]
			offset--
		}
	}

	if vmoptionsContentString[len(vmoptionsContentString)-1:] != "\n" {
		vmoptionsContentString += "\n"
	}

	return vmoptionsContentString, agents
}

func getKeyIndexBySlug(slug string) (int, error) {
	for i, v := range KeyListSlugIndexed {
		if v == slug {
			return i, nil
		}
	}
	return 0, errors.New("there is no key for " + slug)
}

func fileHash(path string) (uint32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	return binaryHash(data), nil
}

func binaryHash(data []byte) uint32 {
	return adler32.Checksum(data)
}

func checkAgentExists(agentPaths []string) bool {
	agentData := getResource(agentName)
	agentHash := binaryHash(agentData)

	for _, agentPath := range agentPaths {
		hash, err := fileHash(agentPath)
		if err != nil {
			continue
		}

		if hash == agentHash {
			return true
		}
	}

	return false
}
