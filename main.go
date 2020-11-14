package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(srcPath string) (err error) {
	destPath := strings.TrimSuffix(srcPath, ".txt") + ".html"

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, data, 0644)
	return err
}

func main() {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				name := info.Name()
				if strings.HasSuffix(name, ".txt") {
					err := copyFile(path)
					if err != nil {
						return err
					}
				}
			}
			return err
		})

	if err != nil {
		log.Println(err)
	}
}
