package main

import (
	"os"
	"path/filepath"
	"strings"
)

func walkThroughTextFiles(
	walker func(string, os.FileInfo) error,
) error {
	return filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				name := info.Name()
				if strings.HasSuffix(name, ".txt") {
					err := walker(path, info)
					if err != nil {
						return err
					}
				}
			}
			return err
		})
}

func expandHomeDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		return filepath.Join(dirname, path[2:])
	}
	return path
}
