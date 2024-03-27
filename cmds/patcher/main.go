package main

import (
	"flag"
	"log"

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
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetPrefix("\x1b[32m[Pyon]\x1b[0m ")

	patcher.PatchDiscord(&ipaFile, &iconsFile)
}
