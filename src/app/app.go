package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/inancgumus/screen"
	"log"
	"os"
	"runtime"
	"strconv"
)

//go:embed resources_enc
var resources embed.FS

var osName string

var stdin *bufio.Reader

func init() {
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
}

func main() {
	menu_loop()
}

func menu_loop() {
	const ITEM_SHOWINFO = 0
	const ITEM_PATCH = 1
	const ITEM_EXIT = 2

	items := map[byte]string{
		ITEM_SHOWINFO: "Show info",
		ITEM_PATCH:    "Patch (default flow)",
		ITEM_EXIT:     "Exit",
	}

	func() {
		var selected_item int
		for {
			screen.Clear()
			screen.MoveTopLeft()

			for i := 0; i < len(items); i++ {
				fmt.Printf("%d. %s\n", i, items[byte(i)])
			}

			var inbuf []byte
			inbuf, _, _ = stdin.ReadLine()
			selected_item, _ = strconv.Atoi(string(inbuf))

			if selected_item == ITEM_EXIT {
				break
			}

			switch selected_item {
			case ITEM_SHOWINFO:
				item_show_info()
				break
			case ITEM_PATCH:
				item_patch()
				break
			}

			delay()
		}
	}()
}
