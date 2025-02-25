package main

func item_show_info_get_text() string {
	info := getOgurchiki() + "\n\n" + "Music: " + getKapusta() + "\n\n" + getPerchik()
	return info
}
