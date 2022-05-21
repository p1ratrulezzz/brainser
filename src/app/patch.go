package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const agentName = "agent.jar"

func doPatch(vmoptionsPath string, destinationPath string, keyIndex int) []string {
	var errorMessages []string

	destinationDir := destinationPath
	if destinationDir == "" {
		destinationDir = filepath.Dir(vmoptionsPath)
	} else if _, err := os.Stat(destinationDir); err != nil {
		fmt.Println("Folder " + destinationDir + " doesn't exists. Will create it")
		os.Mkdir(destinationDir, 0755)
	}

	jarname := filepath.Join(destinationDir, agentName)

	vmoptionsName := filepath.Base(vmoptionsPath)
	vmoptionsNewPath := filepath.Join(destinationDir, vmoptionsName)

	var vmoptionsContent []byte
	if _, err := os.Stat(vmoptionsNewPath); errors.Is(err, os.ErrNotExist) {
		vmoptionsContent, _ = os.ReadFile(vmoptionsPath)
	} else {
		vmoptionsContent, _ = os.ReadFile(vmoptionsNewPath)
	}

	vmoptionsContentString, agents := cleanupVmoptions(vmoptionsContent)
	if !checkAgentExists(agents) {
		vmoptionsContentString += "-javaagent:" + jarname
	} else {
		errorMessages = append(errorMessages, "This product is already patched")
		return errorMessages
	}

	if _, err := os.Stat(jarname); err == nil {
		fmt.Printf("File %s already exists, but will be overwritten\n", jarname)
	}

	jarfileContent := getResource(agentName)
	err := os.WriteFile(jarname, jarfileContent, 0644)
	if err != nil {
		errorMessages = append(errorMessages, fmt.Sprintf("File %s can't be written", jarname))
		return errorMessages
	}

	err = os.WriteFile(vmoptionsNewPath, []byte(vmoptionsContentString), 0644)
	if err != nil {
		errorMessages = append(errorMessages, fmt.Sprintf("Writing error. Error: "+err.Error()))
		return errorMessages
	}

	keyPath := filepath.Join(destinationDir, KeyListSlugIndexed[keyIndex]+".key")
	fpKey, err := os.Create(keyPath)
	keyContent := getResource("universal.key")
	fpKey.Write(keyContent)
	fpKey.Close()

	fmt.Println("Patched successfully!")
	return errorMessages
}
