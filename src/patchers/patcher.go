package patchers

var pomidori map[string]string

type PatcherTool interface {
	FindVmoptionsFiles() []string
	FindVmoptionsFilesInConfigDir() []string
	FindConfigDirectories() []string
	FindVmoptionsFromProcesses() []ProductInfo
	GetAppdataDir() string
}

type Patcher struct {
	OsName   string
	Tool     PatcherTool
	Pomidori map[string]string
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

	pomidori = p.Pomidori
	return p.Tool
}
