package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/inancgumus/screen"
	"log"
	"os"
	"reflect"
	"runtime"
	"sort"
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
	// menu_loop()
	gui()
}

func menu_loop() {
	//item_patch_all()
	// @fixme: Remove
	//panic("remove this")
	const ITEM_SHOWINFO = 0
	const ITEM_PATCH = 1
	const ITEM_PATCH_ALL = 2
	const ITEM_EXIT = 3

	items := map[byte]string{
		ITEM_SHOWINFO: "Show info",
		ITEM_PATCH:    "Patch (default flow)",
		// ITEM_PATCH_ALL: "Smart (but not clever) patch everything",
		ITEM_EXIT: "Exit",
	}

	itemsKeys := reflect.ValueOf(items).MapKeys()
	sort.SliceStable(itemsKeys, func(i, j int) bool {
		return i < j
	})

	func() {
		var selected_item int
		for {
			screen.Clear()
			screen.MoveTopLeft()

			for _, i := range itemsKeys {
				fmt.Printf("%d. %s\n", i, items[byte(i.Uint())])
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
			case ITEM_PATCH_ALL:
				item_patch_all()
				break
			default:
				fmt.Println("Incorrect choice")
				break
			}

			delay()
		}
	}()
}
