package patchers

import (
	"github.com/mattn/go-shellwords"
	"github.com/shirou/gopsutil/process"
	"os"
	"path/filepath"
	"strings"
)

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

func (p *PatcherToolAbstract) GetAppdataDir() string {
	appdataDir, _ := os.UserConfigDir()

	return appdataDir
}

func (p *PatcherToolAbstract) FindVmoptionsFromProcesses() []ProductInfo {
	var infos []ProductInfo

	pids, _ := process.Pids()
	for _, pid := range pids {
		proc, _ := process.NewProcess(pid)
		cmdline, err := proc.Cmdline()
		if err != nil {
			continue
		}

		if strings.Index(cmdline, "-Djb.vmOptionsFile=") != -1 &&
			strings.Index(cmdline, "-Didea.platform.prefix=") != -1 &&
			strings.Index(cmdline, ".vmoptions") != -1 {
			parsed, _ := shellwords.Parse(cmdline)
			info := ProductInfo{}
			for _, line := range parsed {
				ruins := strings.Split(line, "=")
				ruins2 := strings.Split(line, ":")
				if len(ruins) > 1 {
					switch ruins[0] {
					case "-Djb.vmOptionsFile":
						info.VmoptionsSourcePath = ruins[1]
						break
					case "-Didea.platform.prefix":
						info.ProductName = ruins[1]
						info.ProductSlug = strings.ToLower(ruins[1])
						break
					case "-Didea.paths.selector":
						info.ProductFolder = ruins[1]
						info.VmoptionsDestinationPath = filepath.Join(p.GetAppdataDir(), "JetBrains", ruins[1])
						break
					}
				}

				if len(ruins2) > 1 {
					switch ruins2[0] {
					case "-javaagent":
						info.Agents = append(info.Agents, ruins2[1])
						break
					}
				}
			}

			infos = append(infos, info)
		}
	}

	return infos
}

func (p *PatcherToolAbstract) FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func (p *PatcherToolAbstract) GetExeList() *map[string]string {
	var exeListPtr *map[string]string

	func() {
		if exeListPtr != nil {
			return
		}

		var exeList = map[string]string{
			"appcode":   "AppCode",
			"clion":     "Clion",
			"datagrip":  "Datagrip",
			"dataspell": "DataSpell",
			"goland":    "GoLand",
			"idea":      "Idea",
			"phpstorm":  "PhpStorm",
			"pycharm":   "PyCharm",
			"rider":     "Rider",
			"rubymine":  "RubyMine",
			"webstorm":  "WebStorm",
		}

		exeListPtr = &exeList
	}()

	return exeListPtr
}
