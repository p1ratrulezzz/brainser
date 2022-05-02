package patchers

import "os"

type PatcherToolLinux struct {
	PatcherToolAbstract
}

func (p *PatcherToolLinux) FindVmoptionsFiles() []string {
	homeDir, _ := os.UserHomeDir()
	files := findVmoptionsFiles([]string{"/opt", homeDir + "/.local/share/JetBrains"})

	return files
}
