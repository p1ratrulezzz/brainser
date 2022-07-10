package main

import (
	"fmt"
	"jetbrainser/src/patchers"
	"strings"
)

func item_patch_procs() {
	if globalvarCleanupMode {
		fmt.Println("Attention! Working in cleanup mode. This will not actually patch anything, just cleanup other agents! Disable cleanup mode to patch")
	}

	fmt.Println("Run all products that you want to patch and press enter...")
	stdin.ReadLine()

	fmt.Println("Searching for products...")
	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()

	allProducts := tool.FindVmoptionsFromProcesses()

	if len(allProducts) == 0 {
		fmt.Println("No running products found")
		return
	}

	fmt.Printf("Found %d products\n", len(allProducts))
	fmt.Println("Will patch these products:")
	for _, info := range allProducts {
		fmt.Printf("%s(%s)\n", info.ProductName, info.VmoptionsSourcePath)
	}

	fmt.Println("Press enter to continue...")
	stdin.ReadLine()

	messages := doAutoPatch(tool, allProducts)
	fmt.Println(strings.Join(messages, "\n"))
}

func doAutoPatch(tool patchers.PatcherTool, allProducts []patchers.ProductInfo) []string {
	var messages []string

	for _, info := range allProducts {
		messages = append(messages, fmt.Sprintf("Found product %s\n", info.ProductName))

		keyIndex, err := getKeyIndexBySlug(info.ProductSlug)
		if err != nil {
			messages = append(messages, fmt.Sprintf("There is no key for %s\n", info.ProductName))
			continue
		}

		messages = append(messages, "Patching...")
		errorMessages := doPatch(info.VmoptionsSourcePath, info.VmoptionsDestinationPath, info.AgentDir, keyIndex)
		messages = append(messages, errorMessages...)
	}

	messages = append(messages, "All products patched! Close all your products and run again.")
	return messages
}
