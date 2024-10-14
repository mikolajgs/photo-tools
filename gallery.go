package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mikolajgs/broccli"
)

var templates = template.Must(template.ParseFiles("index.html", "gallery.html"))

type Directory struct {
	Name string
}

type Photo struct {
	ThumbPath string
	PhotoPath string
}

func serveGalleryHandler(c *broccli.CLI) int {
	photos := c.Flag("photos")
	thumbs := c.Flag("thumbs")

	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir(photos))))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir(thumbs))))

	// Serve index and gallery pages
	http.HandleFunc("/", getIndexHandler(photos))
	http.HandleFunc("/gallery/", getGalleryHandler(photos))

	// Start the server
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	return 0
}

func getIndexHandler(photoDir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dirs, err := getSubdirectories(photoDir)
		if err != nil {
			http.Error(w, "Unable to read directories", http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "index.html", dirs)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
	}
}

func getGalleryHandler(photoDir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the subdirectory (album) from the URL path
		subdir := strings.TrimPrefix(r.URL.Path, "/gallery/")
		if subdir == "" {
			http.NotFound(w, r)
			return
		}

		thumbs, photos, err := getPhotosInSubdirectory(subdir, photoDir)
		if err != nil {
			http.Error(w, "Unable to read photos", http.StatusInternalServerError)
			return
		}

		var gallery []Photo
		for i := range thumbs {
			gallery = append(gallery, Photo{ThumbPath: thumbs[i], PhotoPath: photos[i]})
		}

		err = templates.ExecuteTemplate(w, "gallery.html", gallery)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
	}
}

func getSubdirectories(root string) ([]Directory, error) {
	var dirs []Directory
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, Directory{Name: entry.Name()})
		}
	}

	return dirs, nil
}

func getPhotosInSubdirectory(subdir string, photoDir string) ([]string, []string, error) {
	photoPath := filepath.Join(photoDir, subdir)

	thumbs := []string{}
	photos := []string{}

	entries, err := os.ReadDir(photoPath)
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && isImage(entry.Name()) {
			thumbs = append(thumbs, filepath.Join("/thumbnails", subdir, entry.Name()+"_th.jpg"))
			photos = append(photos, filepath.Join("/photos", subdir, entry.Name()))
		}
	}

	return thumbs, photos, nil
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".heif" || ext == ".heic"
}
