package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/muesli/termenv"
)

const (
	promptTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- print "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- print "⇣ " -}} 
  {{- else -}}
    {{- print "  " -}}
  {{- end -}} 

  {{- if eq $.SelectedIndex $i }}
   {{- print "[" (Foreground "32" (Bold "x")) "] " (Selected $choice) "\n" }}
  {{- else }}
    {{- print "[ ] " (Unselected $choice) "\n" }}
  {{- end }}
{{- end}}
`
	resultTemplate = `
		{{- print .Prompt " " (Foreground "32"  (name .FinalChoice)) "\n" -}}
		`
)

type Movie struct {
	TMDBID       string
	IMDBID       string
	Title        string
	EnglishTitle string
	Year         int
	Directors    []string
	Summary      string
}

func main() {
	tmdb, err := NewTMDB(os.Getenv("TMDB_API_KEY"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input := textinput.New("TheMovieDB search term:")
	//input.InitialValue = os.Getenv("USER")
	//input.Placeholder = "Your name cannot be empty"

	searchTerm, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	movies, err := tmdb.Search(searchTerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(movies) == 0 {
		fmt.Println("No results found")
		os.Exit(1)
	}

	blue := termenv.String().Foreground(termenv.ANSI256Color(32)) //nolint:gomnd

	sp := selection.New("Choose a movie!", movies)
	sp.FilterPrompt = "Filter by ID:"
	sp.FilterPlaceholder = "Type to filter"
	sp.PageSize = 3
	sp.LoopCursor = true
	sp.Filter = func(filter string, choice *selection.Choice[Movie]) bool {
		return strings.HasPrefix(choice.Value.TMDBID, filter)
	}
	sp.Template = promptTemplate
	sp.ResultTemplate = resultTemplate
	sp.SelectedChoiceStyle = func(c *selection.Choice[Movie]) string {
		return fmt.Sprintf("%s (%s)", blue.Bold().Styled(c.Value.Title), termenv.String(fmt.Sprintf("%d", c.Value.Year)).Faint().String())
	}
	sp.UnselectedChoiceStyle = func(c *selection.Choice[Movie]) string {
		return fmt.Sprintf("%s (%s)", c.Value.Title, termenv.String(fmt.Sprintf("%d", c.Value.Year)).Faint().String())
	}
	sp.ExtendedTemplateFuncs = map[string]interface{}{
		"name": func(c *selection.Choice[Movie]) string { return c.Value.Title },
	}

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("could not select: %v\n", err)
		os.Exit(1)
	}

	if err := Export(choice); err != nil {
		fmt.Printf("could not export to markdown: %v\n", err)
	}
}
