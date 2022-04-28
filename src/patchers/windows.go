package patchers

import "os"

type PatcherToolWindows struct {
	PatcherTool
}

func (p PatcherToolWindows) FindDirectories() ([]string, []string) {
	configDir, _ := os.UserConfigDir()
	var programfilesDirectories = []string{}
	programfilesDir := os.Getenv("programfiles")
	if programfilesDir != "" {
		programfilesDirectories = append(programfilesDirectories, programfilesDir)
	}

	programfilesDir = os.Getenv("programfiles(x86)")
	if programfilesDir != "" {
		programfilesDirectories = append(programfilesDirectories, programfilesDir)
	}

	_ = programfilesDir
	files := findVmoptionsFiles(programfilesDirectories)
	appdata_dirs := findAppdataDirs(configDir)

	return files, appdata_dirs
}
