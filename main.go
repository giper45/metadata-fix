package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func parseArguments() (string, bool) {
	if len(os.Args) < 2 {
		log.Fatal("Folder path is required")
	}

	folderPath := os.Args[1]
	override := false
	if len(os.Args) > 2 && os.Args[2] == "--override" {
		override = true
	}

	return folderPath, override
}

func checkFolder(folderPath string) {
	fileInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		log.Fatalf("Folder does not exist: %s", folderPath)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("Provided path is not a folder: %s", folderPath)
	}
}

func RunOnFolder(folderPath string, override bool) {
	checkFolder(folderPath)

	images := ReadFolder(folderPath)
	for _, image := range images {
		filename := filepath.Base(image)
		log.Printf("Processing file: %s\n", image)
		if HasMetadataTimeName(filename) {
			modTime := ParseImageTime(filename)
			if modTime.IsZero() {
				log.Printf("Skipping file with unknown date format: %s\n", image)
				continue
			}
			err := ModifyMetadataImage(image, modTime, override)
			if err != nil {
				LogError("Error modifying metadata for file", image, ":", err.Error())
			} else {
				LogOK("Modified metadata for file:", image)
			}
		} else {
			LogWarning("Skipping file with unknown structure: ", image)
		}
		fmt.Println()
		fmt.Println()
	}
}

func main() {
	folderPath, override := parseArguments()
	RunOnFolder(folderPath, override)
}

// w.ShowAndRun()
