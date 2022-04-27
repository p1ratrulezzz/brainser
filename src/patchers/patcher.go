package patchers

import (
	"fmt"
)

type PatcherTool interface {
	FindDirectories()
	PatchDirectory()
}

type Patcher struct {
	OsName string
	Tool   PatcherTool
}

func (p Patcher) GetTool() PatcherTool {
	if p.Tool == nil {
		switch p.OsName {
		case "windows":
			p.Tool = PatcherToolWindows{}
			break
		case "linux":
			p.Tool = PatcherToolLinux{}
			break
		}

	}

	return p.Tool
}

type PatcherToolWindows struct {
	PatcherTool
}

func (p PatcherToolWindows) FindDirectories() {
	fmt.Println("Found some directories for windows")
}

type PatcherToolLinux struct {
	PatcherTool
}

func (p PatcherToolLinux) FindDirectories() {
	fmt.Println("Found some directories for linux")
}
