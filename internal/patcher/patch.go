package patcher

import (
	"errors"
	"log"
	"os"

	"howett.net/plist"
)

func PatchDiscord(discordPath *string, iconsPath *string) {
	log.Println("starting patcher")

	checkFile(discordPath)
	checkFile(iconsPath)

	extractDiscord(discordPath)

	log.Println("renaming Discord to Pyoncord")
	if err := patchName(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Discord renamed")

	log.Println("remove devices whitelist")
	if err := patchDevices(); err != nil {
		log.Fatalln(err)
	}
	log.Println("device whitelist removed")

	log.Println("patch Discord icons")
	extractIcons(iconsPath)
	if err := patchIcon(); err != nil {
		log.Fatalln(err)
	}
	log.Println("icons patched")

	log.Println("showing Discord's document folder in the Files app and Finder/iTunes")
	if err := patchiTunesAndFiles(); err != nil {
		log.Fatalln(err)
	}
	log.Println("patched")

	packDiscord()
	log.Println("cleaning up")
	clearPayload()

	log.Println("done!")
}

// Check if file exists
func checkFile(path *string) {
	_, err := os.Stat(*path)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalln("file not found", *path)
	}
}

// Delete the payload folder
func clearPayload() {
	err := os.RemoveAll("temp")
	if err != nil {
		log.Panicln(err)
	}
}

// Load Discord's plist file
func loadPlist() (map[string]interface{}, error) {
	infoFile, err := os.Open("temp/Payload/Discord.app/Info.plist")
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	decoder := plist.NewDecoder(infoFile)
	if err := decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}

// Save Discord's plist file
func savePlist(info *map[string]interface{}) error {
	infoFile, err := os.OpenFile("temp/Payload/Discord.app/Info.plist", os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	encoder := plist.NewEncoder(infoFile)
	err = encoder.Encode(*info)
	return err
}

// Patch Discord's name
func patchName() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	info["CFBundleName"] = "Enmity"
	info["CFBundleDisplayName"] = "Enmity"

	err = savePlist(&info)
	return err
}

// Remove Discord's device limits
func patchDevices() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	delete(info, "UISupportedDevices")

	err = savePlist(&info)
	return err
}

// Patch the Discord icon to use Enmity's icon
func patchIcon() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}

	icons := info["CFBundleIcons"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})
	icons["CFBundleIconName"] = "EnmityIcon"
	icons["CFBundleIconFiles"] = []string{"EnmityIcon60x60"}

	icons = info["CFBundleIcons~ipad"].(map[string]interface{})["CFBundlePrimaryIcon"].(map[string]interface{})
	icons["CFBundleIconName"] = "EnmityIcon"
	icons["CFBundleIconFiles"] = []string{"EnmityIcon60x60", "EnmityIcon76x76"}

	err = savePlist(&info)
	return err
}

// Show Pyoncord's document folder in Files app and iTunes/Finder
func patchiTunesAndFiles() error {
	info, err := loadPlist()
	if err != nil {
		return err
	}
	info["UISupportsDocumentBrowser"] = true
	info["UIFileSharingEnabled"] = true

	err = savePlist(&info)
	return err
}
