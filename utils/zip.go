package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateZip(dst string, src string) error {
	// Zip the Payload folder
	fmt.Println("Zipping the Payload folder...")
	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(zipFile)
	err = addFilesToZip(zipWriter, src, filepath.Base(src))
	if err != nil {
		return err
	}

	zipWriter.Close()
	zipFile.Close()

	return nil
}

func addFilesToZip(zipWriter *zip.Writer, folderPath string, baseInZip string) error {
	// Get the list of files
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		if file.IsDir() {
			// Add files in sub-folders
			err := addFilesToZip(zipWriter, filePath, filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
		} else {
			// Add file
			err := addFileToZip(zipWriter, filePath, baseInZip)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string, baseInZip string) error {
	// Open the file
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// filename := filepath.Base(filePath)
	filename := filepath.Join(baseInZip, filepath.Base(filePath))
	f, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	// Copy file data to zip writer
	_, err = io.Copy(f, fileToZip)
	if err != nil {
		return err
	}

	return nil
}
