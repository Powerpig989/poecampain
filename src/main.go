package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func exePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path, err := filepath.EvalSymlinks(filepath.Dir(ex))
	if err != nil {
		panic(err)
	}
	return path
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln()
	}

	if err := readAreas(); err != nil {
		log.Fatalln(err)
	}

	guide, err := readGuide()
	if err != nil {
		log.Fatalln(err)
	}

	clientTail, err := startTail(config.Client)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := tea.NewProgram(NewGuideModel(guide, clientTail), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalln(err)
	}
}
