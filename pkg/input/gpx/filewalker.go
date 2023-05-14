package gpx

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FindGPXFiles searches the root filepath for GPX files, and returns a GPX Parser for each one.
func FindGPXFiles(root string) ([]*GPXFile, error) {
	var gpxfiles []*GPXFile
	fileSystem := os.DirFS(root)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".gpx" {
			return nil
		}
		gpxfiles = append(gpxfiles, &GPXFile{
			Filepath: filepath.Join(root, path),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not find gpx files from root '%s': %w", root, err)
	}
	return gpxfiles, nil
}
