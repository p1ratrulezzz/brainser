package patchers

type PatcherTool interface {
	FindVmoptionsFiles() []string
	FindVmoptionsFilesInConfigDir() []string
	FindConfigDirectories() []string
}

type Patcher struct {
	OsName string
	Tool   PatcherTool
}

func (p *Patcher) GetTool() PatcherTool {
	if p.Tool == nil {
		switch p.OsName {
		case "windows":
			p.Tool = &PatcherToolWindows{}
			break
		case "linux":
			p.Tool = &PatcherToolLinux{}
			break
		case "darwin":
			p.Tool = &PatcherToolDarwin{}
			break
		default:
			panic("unknown os")
		}

	}

	return p.Tool
}
