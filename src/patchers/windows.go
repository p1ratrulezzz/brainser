package patchers

import (
	"os"
	"path/filepath"
)

type PatcherToolWindows struct {
	*PatcherToolAbstract
}

func (p *PatcherToolWindows) FindVmoptionsFiles() []string {
	var programfilesDirectories = []string{}
	programfilesDir := os.Getenv("programfiles")
	if programfilesDir != "" {
		programfilesDirectories = append(programfilesDirectories, programfilesDir)
	}

	programfilesDir = os.Getenv("programfiles(x86)")
	if programfilesDir != "" {
		programfilesDirectories = append(programfilesDirectories, programfilesDir)
	}

	configDir, _ := os.UserConfigDir()
	programfilesDirectories = append(programfilesDirectories, filepath.Join(configDir, "..", "Local", "JetBrains", "Toolbox"))

	files := findVmoptionsFiles(programfilesDirectories)

	return files
}
