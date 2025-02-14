package patchers

import (
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type PatcherToolWindows struct {
	*PatcherToolAbstract
}

type ProductInfoJsonLaunch struct {
	os                string
	launcherPath      string
	vmOptionsFilePath string
}

type ProductInfoJson struct {
	BuildNumber       string
	DataDirectoryName string
	Name              string
	ProductCode       string           `json:"productCode"`
	Launch            []map[string]any `json:"launch"`
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

func (p *PatcherToolWindows) FindVmoptionsFromProcesses() []ProductInfo {
	var infos []ProductInfo

	var pomidori = *p.GetPomidori()
	var infosUnuqieMap = make(map[string]byte)

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

		if _, exists := pomidori[exeName]; !exists {
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
		info.VmoptionsSourcePath = filepath.Join(productPath, infoJson.Launch[0]["vmOptionsFilePath"].(string))
		info.VmoptionsDestinationPath = filepath.Join(p.GetAppdataDir(), "JetBrains", infoJson.DataDirectoryName)
		info.ProductFolder = infoJson.DataDirectoryName
		info.ProductName = infoJson.Name
		info.ProductSlug = strings.ToLower(exeName)
		info.BuildNumber = infoJson.BuildNumber

		agentPath := info.VmoptionsDestinationPath

		vmoptionsToolboxPath := filepath.Join(productPath, "../", infoJson.BuildNumber+".vmoptions")
		if p.FileExists(vmoptionsToolboxPath) {
			info.VmoptionsSourcePath = vmoptionsToolboxPath
			info.VmoptionsDestinationPath = ""
			agentPath = vmoptionsToolboxPath
		}

		if !p.isAnsiString(agentPath) {
			info.AgentDir = p.getAlternativeAgentDir()
		}

		if _, err := os.Stat(info.VmoptionsSourcePath); err != nil {
			continue
		}

		if _, ok := infosUnuqieMap[info.VmoptionsDestinationPath]; ok {
			continue
		}

		infosUnuqieMap[info.VmoptionsDestinationPath] = 0
		infos = append(infos, info)
	}

	return infos
}

func (p *PatcherToolWindows) isAnsiString(str string) bool {
	for _, s := range str {
		if utf8.RuneLen(s) > 1 {
			return false
		}
	}

	return true
}

func (p *PatcherToolWindows) getAlternativeAgentDir() string {
	programDataFolder := os.Getenv("ProgramData")

	if len(programDataFolder) == 0 || !p.FileExists(programDataFolder) {
		driveLetter := os.Getenv("systemdrive")
		programDataFolder = filepath.Join(driveLetter+":", "ProgramData")
	}

	programDataFolder = filepath.Join(programDataFolder, "JetBrainser")

	os.Mkdir(programDataFolder, 0755)
	return programDataFolder
}
