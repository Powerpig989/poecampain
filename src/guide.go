package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

var statePath = expandTilde(path.Join("~", ".cache", "poecampain"))

type Step struct {
	Text     string
	NextZone string
}

type Guide struct {
	Steps []Step
	Pos   int
}

func (g *Guide) SetPos(p int) {
	if p < 0 || p >= len(g.Steps) {
		g.Pos = 0
	}
	g.Pos = p
}

func (g *Guide) Next() {
	if g.Pos < len(g.Steps)-1 {
		g.Pos = g.Pos + 1
	}
}

func (g *Guide) Prev() {
	if g.Pos > 0 {
		g.Pos = g.Pos - 1
	}
}

func (g *Guide) Start() {
	g.Pos = 0
}

func (g *Guide) End() {
	g.Pos = len(g.Steps) - 1
}

func (g Guide) Display() string {
	return g.Steps[g.Pos].Text
}

func (g Guide) NextZone() string {
	return g.Steps[g.Pos].NextZone
}

func (g Guide) IsNextZone(zone string) bool {
	return g.Steps[g.Pos].NextZone == zone
}

func NewGuide(steps []Step) Guide {
	return Guide{
		Steps: steps,
		Pos:   0,
	}
}

type RawGuide []string

func (r RawGuide) Parse() *Guide {
	re := regexp.MustCompile(`(\S+)\(([^)]*)\)|\n`)
	steps := make([]Step, 0, len(r))

	for _, entry := range r {
		matches := re.FindAllStringSubmatch(entry, -1)
		if len(matches) == 0 {
			continue
		}

		var nextZone string
		if len(matches) >= 2 {
			areaMatch := matches[len(matches)-2]
			if len(areaMatch) > 2 {
				nextZone = Areas[areaMatch[2]]
			}
		}

		var renderedParts []string
		for i, m := range matches {
			if i == len(matches)-1 {
				continue
			}

			raw, keyword := m[0], m[1]

			if raw == "\n" {
				renderedParts = append(renderedParts, "\n")
				continue
			}

			args := strings.Split(m[2], " ")
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}

			f := keywordOrError(keyword)

			var result string
			switch {
			case m[2] == "":
				result = f()
			case len(args) == 1:
				result = f(args[0])
			case len(args) == 2:
				result = f(args[0], args[1])
			}

			if result != "" {
				renderedParts = append(renderedParts, result)
			}
		}

		var sb strings.Builder
		for i, p := range renderedParts {
			if i > 0 && p != "\n" && !strings.HasSuffix(renderedParts[i-1], "\n") {
				sb.WriteString(" ")
			}
			sb.WriteString(p)
		}
		formattedText := sb.String()

		steps = append(steps, Step{
			NextZone: nextZone,
			Text:     formattedText,
		})
	}

	guide := NewGuide(steps)
	return &guide
}

func readGuide() (*Guide, error) {
	path := filepath.Join(exePath(), "data", "guide.yaml")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var r RawGuide
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode guide: %v", err)
	}

	return r.Parse(), nil
}

func readState() (int, error) {
	data, err := os.ReadFile(expandTilde(statePath))
	if err != nil {
		return 0, err
	}

	pos, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pos, nil
}

func writeState(pos int) error {
	return os.WriteFile(expandTilde(statePath), []byte(strconv.Itoa(pos)), 0644)
}
