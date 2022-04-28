package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const agentName = "JetbrainsIdesCrack_5_3_1_KeepMyLic.jar"

func doPatch(vmoptionsPath string, destinationPath string, keyIndex int) {
	destinationDir := destinationPath
	if destinationDir == "" {
		destinationDir = filepath.Dir(vmoptionsPath)
	}

	jarname := filepath.Join(destinationDir, agentName)
	if _, err := os.Stat(jarname); err == nil {
		fmt.Printf("File %s already exists\n", jarname)
	} else {
		fpJarfile, err := os.Create(jarname)
		if err != nil {
			fmt.Printf("Can't create file %s. Error: %s \n", jarname, err.Error())
			return
		}

		jarfileContent, _ := resources.ReadFile(filepath.Join("resources", agentName))
		fpJarfile.Write(jarfileContent)
		fpJarfile.Close()
	}

	vmoptionsName := filepath.Base(vmoptionsPath)
	vmoptionsNewPath := filepath.Join(destinationDir, vmoptionsName)

	var vmoptionsContent []byte
	if _, err := os.Stat(vmoptionsNewPath); errors.Is(err, os.ErrNotExist) {
		vmoptionsContent, _ = os.ReadFile(vmoptionsPath)
	} else {
		vmoptionsContent, _ = os.ReadFile(vmoptionsNewPath)
	}

	vmoptionsContentString := string(vmoptionsContent)
	offset := 0
	needle := "-javaagent:"
	for pos := 0; offset < len(vmoptionsContentString) && pos != -1; {
		pos = strings.Index(vmoptionsContentString[offset:], needle)
		if pos == -1 {
			continue
		}

		pos += offset

		if vmoptionsContentString[(pos-1):pos] != "#" {
			vmoptionsContentString = vmoptionsContentString[0:pos] + "#" + vmoptionsContentString[pos:]
			offset--
		}

		offset += pos + len(needle)
	}

	if vmoptionsContentString[len(vmoptionsContentString)-1:] != "\n" {
		vmoptionsContentString += "\n"
	}

	vmoptionsContentString += needle + jarname

	err := os.WriteFile(vmoptionsNewPath, []byte(vmoptionsContentString), 0644)
	if err != nil {
		fmt.Println("Writing error. Error: " + err.Error())
		return
	}

	keys, _ := getKeys()
	key := keys[keyIndex]

	keyPath := filepath.Join(destinationDir, key.Name())
	fpKey, _ := os.Create(keyPath)
	keyContent, _ := resources.ReadFile(filepath.Join("resources", "keys", key.Name()))
	fpKey.Write(keyContent)
	fpKey.Close()

	fmt.Println("Patched successfully!")
}
