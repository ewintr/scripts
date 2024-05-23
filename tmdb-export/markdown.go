package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	markdownTemplate = `---
tmdb: {{ .TMDBID }}
emdb: {{ .IMDBID }}
englishTitle: {{ .EnglishTitle }}
title: {{ .Title }}
year: {{ .Year }}
runtime: {{ .Runtime }}
directors: {{ .DirectorsYAML }}
inCollection: no
watchedOn:
rating:
---

# {{ .EnglishTitle }} ({{ .Year }})

Director(s): {{ .Directors }}

{{ .Summary }}

## Comment
`
)

func Export(movie Movie) error {
	tpl, err := template.New("page").Parse(markdownTemplate)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s (%d).md", movie.EnglishTitle, movie.Year)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	runtime := time.Duration(movie.RunTime) * time.Minute
	runtimeStr, _, ok := strings.Cut(runtime.String(), "m")
	if !ok {
		return fmt.Errorf("could not parse runtime format %s", runtime)
	}
	runtimeStr = fmt.Sprintf("%sm", runtimeStr)

	data := struct {
		TMDBID        string
		IMDBID        string
		EnglishTitle  string
		Title         string
		Year          int
		Runtime       string
		DirectorsYAML string
		Directors     string
		Summary       string
	}{
		TMDBID:        movie.TMDBID,
		IMDBID:        movie.IMDBID,
		EnglishTitle:  movie.EnglishTitle,
		Title:         movie.Title,
		Year:          movie.Year,
		Runtime:       runtimeStr,
		DirectorsYAML: strings.Join(movie.Directors, ", "),
		Directors:     fmt.Sprintf("[[%s]]", strings.Join(movie.Directors, "]], [[")),
		Summary:       movie.Summary,
	}

	if err := tpl.Execute(f, data); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
