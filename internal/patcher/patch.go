package patcher

import (
	"errors"
	"log"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"fmt"

	"howett.net/plist"
)

const (
	DEFAULT_IPA_PATH    = "files/Discord.ipa"
	DEFAULT_ICONS_PATH  = "files/icons.zip"
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

	log.Println("patching react navigation elements")
	if err := patchReactNavigationElements(); err != nil {
		log.Fatalln(err)
	}
	log.Println("react navigation patched")

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

func patchReactNavigationElements() error {
	var reactNavigationPath, reactNavigationFullPath string

	err := filepath.Walk("./temp/Payload/Discord.app/assets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), "@react-navigation+elements@") {
			reactNavigationPath = info.Name()
			reactNavigationFullPath = path
		}

		return nil
	})

	if err != nil {
		return err
	}

	if reactNavigationPath == "" {
		return fmt.Errorf("could not find the @react-navigation+elements folder")
	}

	log.Println("Found the react-navigation elements folder:", reactNavigationPath)

	err = os.Rename(reactNavigationFullPath, reactNavigationFullPath + "/../@react-navigation+elements@patched")
	if err != nil {
		return err
	}

	manifestFile, err := os.OpenFile("temp/Payload/Discord.app/manifest.json", os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer manifestFile.Close()

	manifestData, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return err
	}

	manifestString := string(manifestData)
	manifestString = strings.ReplaceAll(manifestString, reactNavigationPath, "@react-navigation+elements@patched")

	if _, err := manifestFile.Seek(0, 0); err != nil {
		return err
	}

	if _, err := manifestFile.WriteString(manifestString); err != nil {
		return err
	}

	if err := manifestFile.Truncate(int64(len(manifestString))); err != nil {
		return err
	}

	return nil
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
