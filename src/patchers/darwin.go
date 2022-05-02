package patchers

import (
	"os"
)

type PatcherToolDarwin struct {
	*PatcherToolAbstract
}

func (p *PatcherToolDarwin) FindVmoptionsFiles() []string {
	homeDir, _ := os.UserHomeDir()
	files := findVmoptionsFiles([]string{homeDir + "/Applications", "/Applications"})

	return files
}
