package patchers

import (
	"os"
)

type PatcherToolDarwin struct {
	*PatcherToolAbstract
}

func (p *PatcherToolDarwin) FindVmoptionsFiles() []string {
	homeDir, _ := os.UserHomeDir()
	configDir := p.GetAppdataDir()
	files := findVmoptionsFiles([]string{homeDir + "/Applications", "/Applications", configDir + "/Jetbrains/Toolbox"})

	return files
}
