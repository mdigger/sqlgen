package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// SQL описывает непосредственно запрос к базе данных.
type SQL struct {
	Query    string     // исходный текст запроса
	position `yaml:"-"` // позиция с описанием в исходном файле
}

// String возвращает строку с оригинальным SQL запросом.
func (s SQL) String() string {
	return s.Query
}

func (s SQL) node() yaml.Node {
	// представляем разное форматирование строки,
	// в зависимости от того, многострочное описание запроса или однострочное
	var style yaml.Style
	if strings.ContainsAny(s.Query, "\n\r\t") {
		style = yaml.LiteralStyle
	}

	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: style,
		Tag:   "!!str",
		Value: s.Query,
	}
}

// MarshalYAML поддерживает интерфейс [yaml.Marshaler].
func (s SQL) MarshalYAML() (any, error) {
	return s.node(), nil // возвращает строку с описанием SQL запроса
}

// UnmarshalYAML реализует интерфейс [yaml.Unmarshaler].
func (s *SQL) UnmarshalYAML(n *yaml.Node) error {
	s.Query = n.Value           // сохраняем значение запроса
	s.position = parseSource(n) // разбираем и сохраняем позицию запроса в исходном файле

	return nil
}
