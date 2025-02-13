package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var agentStrAdditional []string

func doPatch(vmoptionsPath string, destinationPath string, agentDir string, keyIndex int) []string {
	cleanupOnly := globalvarCleanupMode
	var errorMessages []string

	if cleanupOnly {
		errorMessages = append(errorMessages, fmt.Sprintf("Working in cleanup mode ON! Will not patch anything, only clean"))
	}

	destinationDir := destinationPath
	if destinationDir == "" {
		destinationDir = filepath.Dir(vmoptionsPath)
	} else if _, err := os.Stat(destinationDir); err != nil {
		fmt.Println("Folder " + destinationDir + " doesn't exists. Will create it")
		os.Mkdir(destinationDir, 0755)
	}

	if agentDir == "" {
		agentDir = destinationDir
	}

	jarname := filepath.Join(agentDir, getOvoshi())

	vmoptionsName := filepath.Base(vmoptionsPath)
	vmoptionsNewPath := filepath.Join(destinationDir, vmoptionsName)

	var vmoptionsContent []byte
	if _, err := os.Stat(vmoptionsNewPath); errors.Is(err, os.ErrNotExist) {
		vmoptionsContent, _ = os.ReadFile(vmoptionsPath)
	} else {
		vmoptionsContent, _ = os.ReadFile(vmoptionsNewPath)
	}

	vmoptionsContentString, agents := cleanupVmoptions(vmoptionsContent)
	agentsExist := checkAgentExists(agents)
	if !cleanupOnly {
		if !agentsExist {
			vmoptionsContentString += "-javaagent:" + jarname
		} else {
			vmoptionsContentString = string(vmoptionsContent)
			errorMessages = append(errorMessages, "This product is already patched")
		}
	}

	var err error
	if !cleanupOnly {
		// Add additional strings
		for _, agentStrValue := range agentStrAdditional {
			if strings.Index(vmoptionsContentString, agentStrValue) == -1 {
				vmoptionsContentString += "\n" + agentStrValue
				errorMessages = append(errorMessages, "Adding string "+agentStrValue)
			}
		}

		if _, err = os.Stat(jarname); err == nil {
			errorMessages = append(errorMessages, fmt.Sprintf("File %s already exists, but will be overwritten\n", jarname))
		}

		if !agentsExist {
			jarfileContent := getResource("burger")
			err := os.WriteFile(jarname, jarfileContent, 0644)
			if err != nil {
				errorMessages = append(errorMessages, fmt.Sprintf("File %s can't be written", jarname))
				return errorMessages
			}
		}
	}

	err = os.WriteFile(vmoptionsNewPath, []byte(vmoptionsContentString), 0644)
	if err != nil {
		errorMessages = append(errorMessages, fmt.Sprintf("Writing error. Error: "+err.Error()))
		return errorMessages
	}

	if cleanupOnly {
		errorMessages = append(errorMessages, fmt.Sprintf("Vmoptions cleanup finished!"))
		return errorMessages
	}

	if agentsExist {
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
