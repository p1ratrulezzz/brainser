package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const agentName = "agent.jar"

func doPatch(vmoptionsPath string, destinationPath string, keyIndex int) {
	destinationDir := destinationPath
	if destinationDir == "" {
		destinationDir = filepath.Dir(vmoptionsPath)
	}

	jarname := filepath.Join(destinationDir, agentName)
	if _, err := os.Stat(jarname); err == nil {
		fmt.Printf("File %s already exists, but will be overwritten\n", jarname)
	}

	jarfileContent := getResource(agentName)
	err := os.WriteFile(jarname, jarfileContent, 0644)
	if err != nil {
		fmt.Printf("File %s can't be written", jarname)
		return
	}

	vmoptionsName := filepath.Base(vmoptionsPath)
	vmoptionsNewPath := filepath.Join(destinationDir, vmoptionsName)

	var vmoptionsContent []byte
	if _, err := os.Stat(vmoptionsNewPath); errors.Is(err, os.ErrNotExist) {
		vmoptionsContent, _ = os.ReadFile(vmoptionsPath)
	} else {
		vmoptionsContent, _ = os.ReadFile(vmoptionsNewPath)
	}

	vmoptionsContentString := cleanup_vmoptions(vmoptionsContent)
	vmoptionsContentString += "-javaagent:" + jarname

	err = os.WriteFile(vmoptionsNewPath, []byte(vmoptionsContentString), 0644)
	if err != nil {
		fmt.Println("Writing error. Error: " + err.Error())
		return
	}

	keyPath := filepath.Join(destinationDir, KeyList[keyIndex]+".key")
	fpKey, err := os.Create(keyPath)
	keyContent := getResource("universal.key")
	fpKey.Write(keyContent)
	fpKey.Close()

	fmt.Println("Patched successfully!")
}
