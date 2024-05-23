package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const (
	markdownTemplate = `---
tmdb: {{ .TMDBID }}
emdb: {{ .IMDBID }}
englishTitle: {{ .EnglishTitle }}
title: {{ .Title }}
year: {{ .Year }}
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

	data := struct {
		TMDBID        string
		IMDBID        string
		EnglishTitle  string
		Title         string
		Year          int
		DirectorsYAML string
		Directors     string
		Summary       string
	}{
		TMDBID:        movie.TMDBID,
		IMDBID:        movie.IMDBID,
		EnglishTitle:  movie.EnglishTitle,
		Title:         movie.Title,
		Year:          movie.Year,
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
