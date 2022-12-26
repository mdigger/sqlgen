package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse разбирает файл с описанием запросов и возвращает разобранный результат.
func Parse(filename string) (*Queries, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	dec.KnownFields(true)

	var q Queries
	if err := dec.Decode(&q); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	return &q, nil
}
