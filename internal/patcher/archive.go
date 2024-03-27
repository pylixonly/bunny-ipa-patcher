package patcher

import (
	"compress/flate"
	"log"
	"os"

	"github.com/mholt/archiver"
)

// Extract Discord's IPA
func extractDiscord(discordPath *string) {
	log.Println("extracting", *discordPath)
	format := archiver.Zip{}

	merr := os.Mkdir("temp", 0755)
	if merr != nil {
		log.Fatalln(merr)
	}

	err := format.Unarchive(*discordPath, "./temp")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(*discordPath, "extracted")
}

// Extract Enmity's icons
func extractIcons(iconsPath *string) {
	log.Println("extracting", *iconsPath)

	format := archiver.Zip{}

	err := format.Unarchive(*iconsPath, "temp/Payload/Discord.app/")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(*iconsPath, "extracted")
}

// Pack Discord's IPA
func packDiscord() {
	log.Println("packing discord")

	format := archiver.Zip{
		CompressionLevel: flate.BestCompression,
	}
	err := format.Archive([]string{"temp/Payload"}, "Discord.zip")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Rename("Discord.zip", "Enmity.ipa")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Discord packed")
}
