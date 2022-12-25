package config

import (
	"gopkg.in/yaml.v3"
)

// Query описывает запрос к базе данных.
//
// В нарушение всех инструкций YAML, комментарии из него используются в
// качестве комментариев к запросу, а не вводится отдельное поле с его описанием,
// как должно было бы быть по правилам.
type Query struct {
	Name     string     // название
	Comment  Comment    // комментарий
	Type     Type       // тип запроса
	SQL      SQL        // текст с SQL запросом
	Params   Fields     // список входящих параметров запроса
	Out      Fields     // список исходящих параметров ответа
	position `yaml:"-"` // строка и колонка в исходном файле с SQL запросом
}

// UnmarshalYAML реализует интерфейс [yaml.Unmarshaler].
func (q *Query) UnmarshalYAML(n *yaml.Node) error {
	// запоминаем позицию с описанием запроса в исходном файле
	q.position = parseSource(n)

	for i := 1; i < len(n.Content); i += 2 {
		nameNode, valueNode := n.Content[i-1], n.Content[i]

		// заполняем значения конкретных полей
		switch nameNode.Value {
		case "type":
			if err := q.Type.UnmarshalYAML(valueNode); err != nil {
				return NewError(err, valueNode, "parse type")
			}

		case "sql":
			if err := q.SQL.UnmarshalYAML(valueNode); err != nil {
				return NewError(err, valueNode, "parse sql")
			}

		case "params":
			if err := q.Params.UnmarshalYAML(valueNode); err != nil {
				return NewError(err, valueNode, "parse params")
			}

			// добавляем комментарий, который задан на уровне названия
			q.Params.Comment = parseComments(nameNode)

		case "out":
			if err := q.Out.UnmarshalYAML(valueNode); err != nil {
				return NewError(err, valueNode, "parse out params")
			}

			// добавляем комментарий, который задан на уровне названия
			q.Out.Comment = parseComments(nameNode)

		default:
			return NewError(nil, nameNode, "unknown property %q", nameNode.Value)
		}
	}

	// дополнительные проверки по заполненности полей запросов
	switch q.Type {
	case TypeMany, TypeOne:
		// для запросов, которые возвращают данные, должны быть описаны параметры разбора ответа
		if len(q.Out.Fields) == 0 {
			return NewError(nil, n, "parameters for outgoing data are not described")
		}

		// проверяем, что исходящие параметры не определены как ссылки
		for _, field := range q.Out.Fields {
			switch {
			case field.Type[0] == '*' || field.Type[0] == '&':
				return NewError(nil, n, "unsupported field %q type pointer", field.Name)
			}
		}

	default:
		// для запросов, которые не возвращают данные, параметры ответа не должны быть описаны
		if len(q.Out.Fields) != 0 {
			return NewError(nil, n, "unused parameters for data output are set")
		}
	}

	return nil
}
