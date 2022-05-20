package main

import (
	"fmt"
	"jetbrainser/src/patchers"
)

func item_patch_procs() {
	fmt.Println("Run all products that you want to patch and press enter...")
	stdin.ReadLine()

	fmt.Println("Searching for products...")
	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()
	allProducts := tool.FindVmoptionsFromProcesses()

	fmt.Printf("Found %d products\n", len(allProducts))
	fmt.Println("Will patch these products:")
	for _, info := range allProducts {
		fmt.Printf("%s(%s)\n", info.ProductName, info.VmoptionsSourcePath)
	}

	fmt.Println("Press enter to continue...")
	stdin.ReadLine()

	agentHash := binaryHash(getResource(agentName))
	for _, info := range allProducts {
		fmt.Printf("Found product %s\n", info.ProductName)

		keyIndex, err := getKeyIndexBySlug(info.ProductSlug)
		if err != nil {
			fmt.Printf("There is no key for %s\n", info.ProductName)
			continue
		}

		failedMessage := ""
		for _, agentFile := range info.Agents {
			hashcrc, err := fileHash(agentFile)
			if err != nil {
				failedMessage = "Error calculating hash for agent. Skipping"
				break
			}

			if hashcrc == agentHash {
				failedMessage = "This product is already patched"
				break
			}
		}

		if len(failedMessage) > 0 {
			fmt.Println(failedMessage)
			continue
		}

		fmt.Println("Patching...")
		doPatch(info.VmoptionsSourcePath, info.VmoptionsDestinationPath, keyIndex)
	}

	fmt.Println("All products patched! Close all your products and run again.")
	delay()
}
