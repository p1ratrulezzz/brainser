package patchers

import "os"

type PatcherToolLinux struct {
	PatcherTool
}

func (p PatcherToolLinux) FindDirectories() ([]string, []string) {
	homeDir, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()
	files := findVmoptionsFiles([]string{homeDir, "/opt"})
	appdataDirs := findAppdataDirs(configDir)

	return files, appdataDirs
}
