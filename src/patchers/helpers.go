package patchers

import (
	"os"
	"path/filepath"
	"regexp"
)

func findVmoptionsFiles(paths []string) []string {
	var all_files []string

	for _, root := range paths {
		var files []string
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".vmoptions" {
				files = append(files, path)
			}

			return nil
		})
		if err == nil {
			all_files = append(all_files, files...)
		}
	}

	return all_files
}

func findAppdataDirs(root string) []string {
	var searchPrefixes = map[string][]string{
		"phpstorm": {`(?i)\.?phpstorm[0-9]{4}[0-9\.]+`},
		"idea":     {`(?i)\.?idea[0-9]{4}[0-9\.]+`, `(?mi)\.intellijidea[0-9]{4}[0-9\.]+`},
		"goland":   {`(?i)\.?goland[0-9]{4}[0-9\.]+`},
		"pycharm":  {`(?i)\.?pycharm[0-9]{4}[0-9\.]+`},
		"datagrip":  {`(?i)\.?datagrip[0-9]{4}[0-9\.]+`},
		"rider":  {`(?i)\.?rider[0-9]{4}[0-9\.]+`},
		"rubymine":  {`(?i)\.?rubymine[0-9]{4}[0-9\.]+`},
		"webstorm":  {`(?i)\.?webstorm[0-9]{4}[0-9\.]+`},
		"clion":  {`(?i)\.?clion[0-9]{4}[0-9\.]+`},
	}

	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		filename := info.Name()
		if info.IsDir() {
			for _, patterns := range searchPrefixes {
				for _, pattern := range patterns {
					match, _ := regexp.MatchString(pattern, filename)
					if match {
						files = append(files, path)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return []string{}
	}

	return files
}

func findLinuxAppdataDirs() []string {
	var files []string
	err := filepath.Walk("/home", func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path) == ".config" {
			if err == nil && info.IsDir() {
				files = append(files, findAppdataDirs(path)...)
			}

		}

		return nil
	})

	if err != nil {
		return []string{}
	}

	return files
}
