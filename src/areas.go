package main

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

var Areas AreasData

type AreasData map[string]string

func (a AreasData) GetName(id string) string {
	name, ok := a[id]
	if !ok {
		return ""
	}
	return name
}

func readAreas() error {
	path := filepath.Join(exePath(), "data", "areas.yaml")

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open areas: %v", err)
	}

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&Areas); err != nil {
		return fmt.Errorf("failed to decode areas: %v", err)
	}

	return nil
}
