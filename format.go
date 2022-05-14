package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func formatImageLink(src []byte, srcPath string) []byte {
	srcPathDir := filepath.Dir(srcPath)
	srcPathBase := strings.TrimSuffix(filepath.Base(srcPath), ".txt")

	return imageRegExp.ReplaceAllFunc(src, func(match []byte) []byte {
		text := string(match)
		imagePath := strings.Trim(text, "<> ")

		if strings.HasPrefix(imagePath, "./_img/") {
			return match
		}

		imagePathBase := filepath.Base(imagePath)
		nextImagePathBase := fmt.Sprintf("%s_%s", srcPathBase, imagePathBase)
		nextImagePathDir := filepath.Join(srcPathDir, "_img")
		nextImagePath := filepath.Join(nextImagePathDir, nextImagePathBase)
		nextImagePathRelative := fmt.Sprintf("./_img/%s", nextImagePathBase)

		err := os.MkdirAll(nextImagePathDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		srcFile, err := os.Open(expandHomeDir(imagePath))
		if err != nil {
			fmt.Println(err)
			return match
		}
		defer srcFile.Close()

		destFile, err := os.Create(nextImagePath)
		if err != nil {
			fmt.Println(err)
			return match
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			fmt.Println(err)
			return match
		}

		convertCmd := exec.Command("gm", "convert", nextImagePath, "-resize", "1000x>", nextImagePath)
		convertOut, err := convertCmd.Output()
		if err != nil {
			fmt.Println(convertOut)
			fmt.Println(err)
			return match
		}
		fmt.Println(convertOut)

		tag := fmt.Sprintf(`<< %s >>`, nextImagePathRelative)
		return []byte(tag)
	})
}

func formatFile(srcPath string) error {
	var err error

	src, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	f, err := os.Create(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(formatImageLink(src, srcPath))
	if err != nil {
		return err
	}

	return f.Sync()
}

func format() {
	var err error

	err = walkThroughTextFiles(
		func(path string, info os.FileInfo) error {
			return formatFile(path)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
