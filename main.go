package main

import (
	"fmt"
	"os"

	"github.com/mikolajgs/broccli"
)

var supportedExtensions = []string{".heif", ".heic", ".jpg", ".jpeg", ".png"}


func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func main() {
	cli := broccli.NewCLI("photo-tool", "Tool for handling photos", "")

	showZipFiles := cli.AddCmd("show-zip-files", "Checks various things", showZipFilesHandler)
	showZipFiles.AddFlag("photos", "p", "", "Path to directory with photos", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired|broccli.IsDirectory)

	createThumbs := cli.AddCmd("create-thumbs", "Generate thumbnails", createThumbsHandler)
	createThumbs.AddFlag("photos", "p", "", "Path to directory with photos", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired|broccli.IsDirectory)
	createThumbs.AddFlag("thumbs", "t", "", "Path to thumbs directory", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired|broccli.IsDirectory)

	serveGallery := cli.AddCmd("serve-gallery", "Starts HTTP server with a photo gallery", serveGalleryHandler)
	serveGallery.AddFlag("photos", "p", "", "Path to directory with photos", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired|broccli.IsDirectory)
	serveGallery.AddFlag("thumbs", "t", "", "Path to thumbs directory", broccli.TypePathFile, broccli.IsExistent|broccli.IsRequired|broccli.IsDirectory)

	_ = cli.AddCmd("version", "Shows version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}

	os.Exit(cli.Run())
}
