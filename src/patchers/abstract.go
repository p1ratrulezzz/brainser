package patchers

import (
	"encoding/json"
	"errors"
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

func (p *PatcherToolAbstract) parseProductInfoJson(path string) (*ProductInfoJson, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}
	var dest ProductInfoJson
	err = json.Unmarshal(data, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *PatcherToolAbstract) FindProductInfoJson(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	testJsonFilePath := filepath.Join(path, "product-info.json")
	for !p.FileExists(testJsonFilePath) {
		path, err = filepath.Abs(filepath.Join(path, "../"))
		testJsonFilePath = filepath.Join(path, "product-info.json")
		if err != nil {
			return "", err
		}
	}

	if p.FileExists(testJsonFilePath) {
		return testJsonFilePath, nil
	} else {
		return "", errors.New("no product-info.json file found")
	}
}

func (p *PatcherToolAbstract) GetProductCanonicalNameByCode(productCode string, fallbackName string) string {
	var names = map[string]string{
		"iu": "idea",
	}

	if names[productCode] != "" {
		return names[productCode]
	}

	return fallbackName
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
			strings.Index(cmdline, "-Didea.paths.selector=") != -1 &&
			strings.Index(cmdline, ".vmoptions") != -1 {
			parsed, _ := shellwords.Parse(cmdline)
			productInfoJsonFilePath, err := p.FindProductInfoJson(parsed[0])
			_ = productInfoJsonFilePath
			if err != nil {
				continue
			}

			infoJson, err := p.parseProductInfoJson(productInfoJsonFilePath)
			if err != nil {
				continue
			}

			info := ProductInfo{}
			info.ProductName = infoJson.Name
			info.ProductSlug = strings.ToLower(p.GetProductCanonicalNameByCode(strings.ToLower(infoJson.ProductCode), infoJson.Name))
			info.BuildNumber = infoJson.BuildNumber
			for _, line := range parsed {
				ruins := strings.Split(line, "=")
				ruins2 := strings.Split(line, ":")
				if len(ruins) > 1 {
					switch ruins[0] {
					case "-Djb.vmOptionsFile":
						info.VmoptionsSourcePath = ruins[1]
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

			// Fix toolbox paths
			vmoptionsFilename := filepath.Base(info.VmoptionsSourcePath)
			if vmoptionsFilename == info.BuildNumber+".vmoptions" {
				info.VmoptionsDestinationPath = ""
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
