package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var linkRegExp = regexp.MustCompile(`\[\[([^\[\]]*)\]\]`)

func replaceLinks(src []byte) []byte {
	return linkRegExp.ReplaceAllFunc(src, func(match []byte) []byte {
		text := string(match)
		href := strings.Trim(text, "[] ")
		tag := fmt.Sprintf(`<a href="%s">%s</a>`, href, text)
		return []byte(tag)
	})
}

var defaultTemplate = template.Must(template.New("default").Parse(`<!doctype html>
<head><title>{{.Title}}</title></head><body><pre><code>
{{.Content}}
</code></pre></body>`))

type Template struct {
	Title   string
	Content string
}

func copyFile(srcPath string) (err error) {
	destPath := strings.TrimSuffix(srcPath, ".txt") + ".html"

	src, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}

	err = defaultTemplate.Execute(f, &Template{
		Title:   srcPath,
		Content: string(replaceLinks(src)),
	})
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
