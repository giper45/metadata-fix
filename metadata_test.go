package main

import (
	"io/fs"
	"os"
	"testing"
	"time"
)

func TestMetadataRead(t *testing.T) {
	filePath := "tests/IMG-20171123-WA0012.jpg"
	fileInfo, err := ReadMetadataImage(filePath)

	expectedModTime := time.Date(2025, time.January, 1, 10, 29, 9, 62271471, time.FixedZone("CET", 1*60*60))
	expectedSize := int64(8342)
	expectedMode := fs.FileMode(0644)

	if err != nil {
		t.Error("Error reading metadata")
	}

	if fileInfo.Size() != expectedSize {
		t.Errorf("Expected file size 8342, but got %d", fileInfo.Size())
	}

	if fileInfo.Mode() != expectedMode {
		t.Errorf("Expected file mode -rw-r--r--, but got %v", fileInfo.Mode())
	}

	if !fileInfo.ModTime().Equal(expectedModTime) {
		t.Errorf("Expected modification time %v, but got %v", expectedModTime, fileInfo.ModTime())
	}

}

func TestRenameFile(t *testing.T) {
	testName := "tests/image.jpg"
	newName := RenameFile(testName)
	if newName != "tests/image-modified.jpg" {
		t.Errorf("Expected tests/image-modified.jpg, but got %s", newName)
	}
}

func TestParseImageTime(t *testing.T) {
	tests := []struct {
		filename     string
		expectedTime time.Time
	}{
		{"IMG-20171123-WA0012.jpg", time.Date(2017, time.November, 23, 0, 0, 0, 0, time.UTC)},
		{"VID-20221209-WA0011.mp4", time.Date(2022, time.December, 9, 0, 0, 0, 0, time.UTC)},
		{"IMG-20200101-WA0001.jpg", time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"VID-20181231-WA0002.mp4", time.Date(2018, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"20220821_231703.jpg", time.Date(2022, time.August, 21, 0, 0, 0, 0, time.UTC)},
		{"VID20230123231258.mp4", time.Date(2023, time.January, 23, 0, 0, 0, 0, time.UTC)},
		{"2024-01-01 00.10.25.jpg", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			parsedTime := ParseImageTime(tt.filename)
			if !parsedTime.Equal(tt.expectedTime) {
				t.Errorf("Expected time %v, but got %v", tt.expectedTime, parsedTime)
			}
		})
	}
}

func TestModifyMetadataImage(t *testing.T) {
	filePath := "tests/IMG-20171123-WA0012.jpg"
	modTime := time.Date(2023, time.January, 1, 10, 29, 9, 62271471, time.FixedZone("CET", 1*60*60))
	// Do not replace the file with a new one
	err := ModifyMetadataImage(filePath, modTime, false)

	if err != nil {
		t.Error("Error modifying metadata")
	}
	newFile := RenameFile(filePath)

	fileInfo, err := ReadMetadataImage(newFile)
	if err != nil {
		t.Error("Error reading metadata image")
	}

	if !fileInfo.ModTime().Equal(modTime) {
		t.Errorf("Expected modification time %v, but got %v", modTime, fileInfo.ModTime())
	}
	// Replace the file with a new one
	modTime = time.Date(2023, time.January, 1, 10, 29, 9, 62271471, time.FixedZone("CET", 1*60*60))
	err = ModifyMetadataImage(newFile, modTime, true)
	if err != nil {
		t.Error("Error modifying metadata")
	}
	if !fileInfo.ModTime().Equal(modTime) {
		t.Errorf("Expected modification time %v, but got %v", modTime, fileInfo.ModTime())
	}

	os.Remove("tests/IMG-20171123-WA0012-modified.jpg")
}

func TestHasMetadataTimeName(t *testing.T) {
	tests := []struct {
		filename     string
		expectedTime bool
	}{
		{"IMG-20171123-WA0012.jpg", true},
		{"VID-20221209-WA0011.mp4", true},
		{"IMG-20200101-WA0001.jpg", true},
		{"VID-20181231-WA0002.mp4", true},
		{"20220821_231703.jpg", true},
		{"VID20230123231258.mp4", true},
		{"2024-01-01 00.10.25.jpg", true},
		{"VID-20221209-WA0011.mp4", true},
		{"IMG-20170808-WA0007-modified.jpg", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			hasMetadata := HasMetadataTimeName(tt.filename)
			if hasMetadata != tt.expectedTime {
				t.Errorf("Expected %v, but got %v", tt.expectedTime, hasMetadata)
			}
		})
	}
}
