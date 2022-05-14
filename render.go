package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

var imageRegExp = regexp.MustCompile(`<<([^\[\]]*)>>`)

func replaceImageLinks(src []byte) []byte {
	return imageRegExp.ReplaceAllFunc(src, func(match []byte) []byte {
		text := string(match)
		src := strings.Trim(text, "<> ")
		tag := fmt.Sprintf(`<img src="%s" />`, src)
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

func copyFile(srcPath string, template *template.Template) (err error) {
	destPath := strings.TrimSuffix(srcPath, ".txt") + ".html"

	src, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = template.Execute(f, &Template{
		Title:   srcPath,
		Content: string(replaceImageLinks(replaceLinks(src))),
	})
	return err
}

func readUserTemplate(path string) (*template.Template, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	userTemplate, err := template.New("user").Parse(string(src))
	return userTemplate, err
}

func renderAsHTML(templatePath string) {
	var err error

	template := defaultTemplate

	if len(templatePath) > 0 {
		template, err = readUserTemplate(templatePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = walkThroughTextFiles(
		func(path string, info os.FileInfo) error {
			err := copyFile(path, template)
			return err
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
