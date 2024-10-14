package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikolajgs/broccli"
)

func showZipFilesHandler(c *broccli.CLI) int {
	dir := c.Flag("photos")

	if err := findZipFiles(dir); err != nil {
		fmt.Printf("Error searching for zip files: %v\n", err)
	}

	return 0
}

func findZipFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .zip extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			fmt.Println(path) // Print the full path to the .zip file
		}
		return nil
	})
}
