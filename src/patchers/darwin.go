package patchers

import (
	"os"
)

type PatcherToolDarwin struct {
	PatcherTool
}

func (p PatcherToolDarwin) FindDirectories() ([]string, []string) {
	homeDir, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()
	files := findVmoptionsFiles([]string{homeDir + "/Applications", "/Applications"})
	appdataDirs := findAppdataDirs(configDir)

	return files, appdataDirs
}
