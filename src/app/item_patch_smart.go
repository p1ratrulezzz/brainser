package main

import (
	"fmt"
	"jetbrainser/src/patchers"
	"strings"
)

func item_patch_procs() {
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

	for _, info := range allProducts {
		fmt.Printf("Found product %s\n", info.ProductName)

		keyIndex, err := getKeyIndexBySlug(info.ProductSlug)
		if err != nil {
			fmt.Printf("There is no key for %s\n", info.ProductName)
			continue
		}

		fmt.Println("Patching...")
		errorMessages := doPatch(info.VmoptionsSourcePath, info.VmoptionsDestinationPath, keyIndex)
		fmt.Println(strings.Join(errorMessages, "\n"))
	}

	fmt.Println("All products patched! Close all your products and run again.")
}
