package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ReadFolder(folderPath string) []string {
	var filesWithMetadata []string
	// List of extensions to check for metadata
	metadataExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".heic": true,
		".heif": true,
		".webp": true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".mkv":  true,
		".mp3":  true,
		".wav":  true,
		".flac": true,
		".pdf":  true,
		".docx": true,
	}

	// Read all files and subdirectories in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullPath := filepath.Join(folderPath, file.Name()) // Get full path of the file or directory
		if file.IsDir() {
			// Recursively call ReadFolder for subdirectories
			subfolderImages := ReadFolder(fullPath)
			filesWithMetadata = append(filesWithMetadata, subfolderImages...)
		} else {
			// Check the file extension
			ext := filepath.Ext(file.Name())
			if metadataExtensions[ext] {
				// fmt.Println("Found image:", fullPath) // Print full path of the image
				filesWithMetadata = append(filesWithMetadata, fullPath)
			} else {
				LogWarning("[-] Skipping file with ext with no metadata:", fullPath)
			}
		}
	}
	return filesWithMetadata
}

func ReadMetadataImage(filePath string) (fileInfo os.FileInfo, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file metadata
	fileInfo, err = file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func ModifyMetadataImage(filePath string, modTime time.Time, override bool) (err error) {
	if override {
		err = os.Chtimes(filePath, modTime, modTime)
		return err

	} else {
		newFilePath := RenameFile(filePath)
		input, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		err = os.WriteFile(newFilePath, input, 0644)
		if err != nil {
			return err
		} else {
			return os.Chtimes(newFilePath, modTime, modTime)
		}
	}
}

func RenameFile(filePath string) string {
	newFilePathext := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filePath, newFilePathext)
	newFilePath := fileName + "-modified" + newFilePathext
	return newFilePath
}

func ParseImageTime(filename string) (expectedTime time.Time) {
	parsers := []struct {
		pattern string
		parser  func(string) (time.Time, error)
	}{
		{
			pattern: `^IMG-(\d{8})-WA\d+(-[a-zA-Z0-9]+)?\.\w+$`,
			parser:  parseDateYYYYMMDD,
		},
		{
			pattern: `^VID-(\d{8})-WA\d+(-[a-zA-Z0-9]+)?\.\w+$`,
			parser:  parseDateYYYYMMDD,
		},
		{
			pattern: `^(\d{8})_\d{6}(-[a-zA-Z0-9]+)?\.\w+$`,
			parser:  parseDateYYYYMMDD,
		},
		{
			pattern: `^VID(\d{8})\d{6}(-[a-zA-Z0-9]+)?\.\w+$`,
			parser:  parseDateYYYYMMDD,
		},
		{
			pattern: `^(\d{4}-\d{2}-\d{2}) \d{2}\.\d{2}\.\d{2}(-[a-zA-Z0-9]+)?\.\w+$`,
			parser:  parseDateYYYY_MM_DD,
		},
	}

	for _, p := range parsers {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindStringSubmatch(filename)
		if matches != nil {
			datePart := matches[1]
			var err error
			expectedTime, err = p.parser(datePart)
			if err != nil {
				log.Printf("Error parsing date from filename '%s': %v\n", filename, err)
				return time.Time{}
			}
			return expectedTime
		}
	}

	LogWarning("Warning: Unknown filename structure", filename)
	return time.Time{} // Return zero value for unknown structures
}

func HasMetadataTimeName(filename string) bool {
	parsers := []struct {
		pattern string
	}{
		{pattern: `^IMG-(\d{8})-WA\d+(-[a-zA-Z0-9]+)?\.\w+$`},
		{pattern: `^VID-(\d{8})-WA\d+(-[a-zA-Z0-9]+)?\.\w+$`},
		{pattern: `^(\d{8})_\d{6}(-[a-zA-Z0-9]+)?\.\w+$`},
		{pattern: `^VID(\d{8})\d{6}(-[a-zA-Z0-9]+)?\.\w+$`},
		{pattern: `^(\d{4}-\d{2}-\d{2}) \d{2}\.\d{2}\.\d{2}(-[a-zA-Z0-9]+)?\.\w+$`},
	}

	for _, p := range parsers {
		re := regexp.MustCompile(p.pattern)
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}

func parseDateYYYYMMDD(datePart string) (time.Time, error) {
	return time.Parse("20060102", datePart)
}

func parseDateYYYY_MM_DD(datePart string) (time.Time, error) {
	return time.Parse("2006-01-02", datePart)
}
