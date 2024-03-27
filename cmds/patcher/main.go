package main

import (
	"flag"

	"io.github.pyoncord/patcher/internal/patcher"
)

var (
	ipaFile   string
	iconsFile string
)

func init() {
	flag.StringVar(&ipaFile, "d", patcher.DEFAULT_IPA_PATH, "Path for Discord.ipa")
	flag.StringVar(&iconsFile, "i", patcher.DEFAULT_ICONS_PATH, "Path for icons.zip")

	flag.Parse()
}

func main() {
	patcher.PatchDiscord(&ipaFile, &iconsFile)
}
