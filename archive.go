package main

import (
	"compress/flate"
	"log"
	"os"

	"github.com/mholt/archiver"
)

// Extract Discord's IPA
func extractDiscord(discordPath *string) {
	log.Println("Extracting", *discordPath)
	format := archiver.Zip{}

	if _, err := os.Stat(".temp"); !os.IsNotExist(err) {
		os.RemoveAll(".temp")
	}

	merr := os.Mkdir(".temp", 0755)
	if merr != nil {
		log.Fatalln(merr)
	}

	err := format.Unarchive(*discordPath, "./.temp")
	if err != nil {
		log.Fatalln(err)
	}
}

// Extract Pyoncord's icons
func extractIcons(iconsPath *string) {
	log.Println("Extracting", *iconsPath)

	format := archiver.Zip{}

	err := format.Unarchive(*iconsPath, ".temp/Payload/Discord.app/")
	if err != nil {
		log.Fatalln(err)
	}
}

// Pack Discord's IPA
func packDiscord() {
	log.Println("Packing IPA")

	format := archiver.Zip{
		CompressionLevel: flate.BestCompression,
	}
	err := format.Archive([]string{".temp/Payload"}, "Discord.zip")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Rename("Discord.zip", "Pyoncord.ipa")
	if err != nil {
		log.Fatalln(err)
	}
}
