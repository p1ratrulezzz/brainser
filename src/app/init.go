package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"
)

var osName string

//go:embed resources_enc
var resources embed.FS

var stdin *bufio.Reader

var globalvarCleanupMode bool

func init() {
	globalvarCleanupMode = false
	osName = runtime.GOOS
	os_supported := map[string]bool{
		"windows": true,
		"linux":   true,
		"darwin":  true,
	}

	if !os_supported[osName] {
		log.Fatal(fmt.Sprintf("this os \"%s\" is not supported (yet)", osName))
	}

	stdin = bufio.NewReader(os.Stdin)
	for slug, name := range KeyList {
		KeyListSlugIndexed = append(KeyListSlugIndexed, slug)
		KeyListNameIndexed = append(KeyListNameIndexed, name)
	}

	// Init additional strings
	agentStrAdditional = append(agentStrAdditional, "--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED")
	agentStrAdditional = append(agentStrAdditional, "--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED")
}
