package patchers

import "os"

type PatcherToolAbstract struct {
	PatcherTool
}

func (p *PatcherToolAbstract) FindConfigDirectories() []string {
	configDir, _ := os.UserConfigDir()
	appdataDirs := findAppdataDirs(configDir)

	return appdataDirs
}

func (p *PatcherToolAbstract) FindVmoptionsFilesInConfigDir() []string {
	configDir, _ := os.UserConfigDir()
	files := findVmoptionsFiles([]string{configDir})

	return files
}
