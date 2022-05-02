package main

import (
	"jetbrainser/src/patchers"
	"path/filepath"
)

func item_patch_all() {
	patcher := patchers.Patcher{osName, nil}
	tool := patcher.GetTool()
	allProducts := tool.FindVmoptionsFiles()
	for _, path := range allProducts {
		filepath.Dir(path)
	}

	_ = tool
}
