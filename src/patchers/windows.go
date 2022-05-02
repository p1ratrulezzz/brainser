package patchers

import "os"

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

	files := findVmoptionsFiles(programfilesDirectories)

	return files
}
