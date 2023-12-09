package polarsteps

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FindPolarstepsFiles searches the root filepath for locations.json files, and returns a polarstep Parser for each one.
func FindPolarstepsFiles(root string) ([]*PolarstepLocationFile, error) {
	var plfs []*PolarstepLocationFile
	fileSystem := os.DirFS(root)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "locations.json" {
			return nil
		}
		plfs = append(plfs, &PolarstepLocationFile{
			Filepath: filepath.Join(root, path),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not find polarstep location files from root '%s': %w", root, err)
	}
	return plfs, nil
}
