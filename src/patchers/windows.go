package patchers

import (
	"encoding/json"
	"github.com/shirou/gopsutil/process"
	"os"
	"path/filepath"
	"strings"
)

type PatcherToolWindows struct {
	*PatcherToolAbstract
}

type ProductInfoJsonLaunch struct {
	VmOptionsFilePath string
}

type ProductInfoJson struct {
	BuildNumber       string
	DataDirectoryName string
	Name              string
	Launch            []ProductInfoJsonLaunch
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

	configDir, _ := os.UserConfigDir()
	programfilesDirectories = append(programfilesDirectories, filepath.Join(configDir, "..", "Local", "JetBrains", "Toolbox"))

	files := findVmoptionsFiles(programfilesDirectories)

	return files
}

func (p *PatcherToolWindows) parseProductInfoJson(path string) (*ProductInfoJson, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}
	var dest ProductInfoJson
	err = json.Unmarshal(data, &dest)
	if err != nil {
		return nil, err
	}

	_ = err

	return &dest, nil
}

func (p *PatcherToolWindows) FindVmoptionsFromProcesses() []ProductInfo {
	var infos []ProductInfo

	var exeList = *p.GetExeList()

	pids, _ := process.Pids()
	for _, pid := range pids {
		proc, _ := process.NewProcess(pid)
		exeNameConst, err := proc.Exe()
		if err != nil {
			continue
		}

		exeName := filepath.Base(exeNameConst)
		ext := filepath.Ext(exeName)
		exeName = exeName[0:(len(exeName) - len(ext))]
		if exeName[(len(exeName)-2):] == "64" {
			exeName = exeName[0:(len(exeName) - 2)]
		}

		if _, exists := exeList[exeName]; !exists {
			continue
		}

		productPath, err := filepath.Abs(filepath.Join(exeNameConst, "../../"))
		if err != nil {
			continue
		}

		productInfoFilePath := filepath.Join(productPath, "product-info.json")
		if _, err := os.Stat(productInfoFilePath); err != nil {
			continue
		}

		infoJson, err := p.parseProductInfoJson(productInfoFilePath)
		if err != nil || len(infoJson.Launch) == 0 {
			continue
		}

		var info ProductInfo
		info.VmoptionsSourcePath = filepath.Join(productPath, infoJson.Launch[0].VmOptionsFilePath)
		info.VmoptionsDestinationPath = filepath.Join(p.GetAppdataDir(), "JetBrains", infoJson.DataDirectoryName)
		info.ProductFolder = infoJson.DataDirectoryName
		info.ProductName = infoJson.Name
		info.ProductSlug = strings.ToLower(exeName)

		vmoptionsToolboxPath := filepath.Join(productPath, "../", infoJson.BuildNumber+".vmoptions")
		if p.FileExists(vmoptionsToolboxPath) {
			info.VmoptionsSourcePath = vmoptionsToolboxPath
			info.VmoptionsDestinationPath = ""
		}

		if _, err := os.Stat(info.VmoptionsSourcePath); err != nil {
			continue
		}

		infos = append(infos, info)
	}

	return infos
}
