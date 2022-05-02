package patchers

import (
	"os"
)

type PatcherToolDarwin struct {
	*PatcherToolAbstract
}

func (p *PatcherToolDarwin) FindVmoptionsFiles() []string {
	homeDir, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()
	files := findVmoptionsFiles([]string{homeDir + "/Applications", "/Applications", configDir + "/Jetbrains/Toolbox"})

	return files
}
