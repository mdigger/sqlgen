package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// position описывает ссылку на строку и позицию в исходном описании YAML.
type position [2]int

// Source возвращает строковое представление номера строки и позиции в исходном файле.
func (s position) Source() string {
	return fmt.Sprintf("%d:%d", s[0], s[1])
}

// parseSource
func parseSource(n *yaml.Node) position {
	return position{n.Line, n.Column}
}
