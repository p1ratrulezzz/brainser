package main

import (
	"bufio"
	"embed"
	"fmt"
	"jetbrainser/src/patchers"
	"log"
	"os"
	"runtime"
)

//go:embed resources
var resources embed.FS

var os_name string

func init() {
	os_name = runtime.GOOS
	os_supported := map[string]bool{
		"windows": true,
		"linux":   true,
		"darwin":  true,
	}

	if !os_supported[os_name] {
		log.Fatal(fmt.Sprintf("this os \"%s\" is not supported (yet)", os_name))
	}
}

func delay() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	// data, _ := app.resources.ReadFile("app.resources/JetbrainsIdesCrack_5_3_1_KeepMyLic.jar")
	// print(string(data))
	print("OK")

	patcher := patchers.Patcher{os_name, nil}

	patcher.GetTool().FindDirectories()

	delay()
}
