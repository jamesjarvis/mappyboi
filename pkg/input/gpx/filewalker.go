package gpx

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindGPXFiles(root string) ([]*GPXFile, error) {
	var gpxfiles []*GPXFile
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".gpx" {
			return nil
		}
		gpxfiles = append(gpxfiles, &GPXFile{
			Filepath: path,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not find gpx files from root '%s': %w", root, err)
	}
	return gpxfiles, nil
}