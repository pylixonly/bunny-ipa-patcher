package main

import (
	"flag"
	"log"
)

var (
	ipaFile   string
	iconsFile string
)

func init() {
	flag.StringVar(&ipaFile, "d", "files/Discord.ipa", "Path for Discord.ipa")
	flag.StringVar(&iconsFile, "i", "files/ipa-icons.zip", "Path for icons.zip")

	flag.Parse()
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetPrefix("\x1b[32m[PyonPatcher]\x1b[0m ")

	PatchDiscord(&ipaFile, &iconsFile)
}
