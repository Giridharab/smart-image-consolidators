package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FindDockerfiles returns paths of all Dockerfiles in the given directory
func FindDockerfiles(dir string) ([]string, error) {
	var dockerfiles []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasPrefix(d.Name(), "Dockerfile") || strings.HasSuffix(d.Name(), ".Dockerfile")) {
			dockerfiles = append(dockerfiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dockerfiles, nil
}

// GetBaseImage extracts the base image from a Dockerfile
func GetBaseImage(dockerfile string) (string, error) {
	file, err := os.Open(dockerfile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 4096)
	n, _ := file.Read(buf)
	content := string(buf[:n])

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "FROM") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("no base image found in %s", dockerfile)
}