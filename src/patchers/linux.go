package patchers

type PatcherToolLinux struct {
	PatcherTool
}

func (p PatcherToolLinux) FindDirectories() ([]string, []string) {
	files := findVmoptionsFiles([]string{"/home", "/opt"})
	appdataDirs := findLinuxAppdataDirs()

	return files, appdataDirs
}
