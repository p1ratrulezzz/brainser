package patchers

import (
	"github.com/shirou/gopsutil/process"
	"howett.net/plist"
	"os"
	"path/filepath"
	"regexp"
)

type PatcherToolDarwin struct {
	*PatcherToolAbstract
}

type ProductInfoPlist struct {
	JVMOptions map[string]interface{} `plist:"JVMOptions"`
}

func (p *PatcherToolDarwin) FindVmoptionsFiles() []string {
	homeDir, _ := os.UserHomeDir()
	configDir := p.GetAppdataDir()
	files := findVmoptionsFiles([]string{homeDir + "/Applications", "/Applications", configDir + "/Jetbrains/Toolbox"})

	return files
}

func (p *PatcherToolDarwin) parsePlist(path string) (*ProductInfoPlist, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0744)
	if err != nil {
		return nil, err
	}

	decoder := plist.NewDecoder(file)
	var infoPlist ProductInfoPlist
	err = decoder.Decode(&infoPlist)
	if err != nil {
		return nil, err
	}

	return &infoPlist, nil
}

func (p *PatcherToolDarwin) FindVmoptionsFromProcesses() []ProductInfo {
	var infos []ProductInfo

	var exeList = *p.GetExeList()
	reList := make(map[string]*regexp.Regexp)
	for product, _ := range exeList {
		re := regexp.MustCompile(`(?i)\.app/Contents/MacOS/` + product + "$")
		reList[product] = re
	}

	pids, _ := process.Pids()
	for _, pid := range pids {
		proc, _ := process.NewProcess(pid)
		exeNameConst, err := proc.Exe()
		if err != nil {
			continue
		}

		for product, re := range reList {
			if re.MatchString(exeNameConst) {
				productPath, _ := filepath.Abs(filepath.Join(exeNameConst, "../../"))
				appFolderName, _ := filepath.Abs(filepath.Join(productPath, "../"))
				toolboxVmoptionsPath := appFolderName + ".vmoptions"
				plistPath := filepath.Join(productPath, "Info.plist")
				plistInfo, err := p.parsePlist(plistPath)
				if err != nil {
					break
				}

				properties := plistInfo.JVMOptions["Properties"].(map[string]interface{})
				pathSelector := properties["idea.paths.selector"].(string)

				var info ProductInfo
				info.ProductSlug = product
				info.ProductName = exeList[product]
				// @todo: Parse product-info.json file and fill build number
				info.ProductFolder = productPath
				info.VmoptionsSourcePath = filepath.Join(productPath, "bin", product+".vmoptions")
				info.VmoptionsDestinationPath = filepath.Join(p.GetAppdataDir(), "JetBrains", pathSelector)

				if !p.FileExists(info.VmoptionsSourcePath) {
					break
				}

				if p.FileExists(toolboxVmoptionsPath) {
					info.VmoptionsSourcePath = toolboxVmoptionsPath
					info.VmoptionsDestinationPath = ""
				}

				infos = append(infos, info)
				break
			}
		}
	}

	return infos
}
