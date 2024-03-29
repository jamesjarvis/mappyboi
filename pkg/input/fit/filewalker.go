package fit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FindFitFiles searches the root filepath for .fit files, and returns a Fit Parser for each one.
func FindFitFiles(root string) ([]*FitFile, error) {
	var fitfiles []*FitFile
	fileSystem := os.DirFS(root)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".fit" {
			fitfiles = append(fitfiles, &FitFile{
				Filepath: filepath.Join(root, path),
			})
			return nil
		}
		if filepath.Ext(path) == ".gz" && strings.Contains(d.Name(), ".fit") {
			fitfiles = append(fitfiles, &FitFile{
				Filepath:   filepath.Join(root, path),
				Compressed: true,
			})
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not find fit files from root '%s': %w", root, err)
	}
	return fitfiles, nil
}
