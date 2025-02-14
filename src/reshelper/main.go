package main

import (
	"fmt"
	"io/fs"
	"jetbrainser/src/cryptor"
	"os"
	"path/filepath"
)

func main() {
	cryptor.SetSauce(os.Getenv("BRAINSER_SAUCE"))
	cryptor.SetSalt(os.Getenv("BRAINSER_SALT"))
	root := "src/app/resources"
	destpath := "src/app/resources_enc"
	var encrypted []byte
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		relpath, _ := filepath.Rel(root, path)
		newpath := filepath.Join(destpath, relpath)
		if info.IsDir() {
			os.Mkdir(newpath, 0755)
		} else {
			newpath += ".enc"
			filedata, _ := os.ReadFile(path)
			if info.Name() == "sauce" || info.Name() == "salt" {
				encrypted = filedata
			} else {
				encrypted = cryptor.Encrypt(filedata)
			}

			os.WriteFile(newpath, encrypted, 0644)
		}

		_ = newpath

		filepath.ToSlash(path)
		return nil
	})

	encrypted, _ = os.ReadFile(destpath + "/check.enc")
	plainText := cryptor.Decrypt(encrypted)

	fmt.Println("Text data " + string(plainText))

	fmt.Println("Encrypted all data")
}
