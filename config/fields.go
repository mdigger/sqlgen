package config

import (
	"gopkg.in/yaml.v3"
)

// Field описывает определение поля запроса.
type Field struct {
	Name     string     // название
	Type     string     // идентификатор типа данных
	Comment  Comment    // комментарий
	position `yaml:"-"` // позиция в исходном файле
}

// Fields описывает список полей запроса.
type Fields struct {
	Comment  Comment        // комментарий
	Fields   []Field        // список полей
	index    map[string]int // индекс с идентификаторами полей
	Anchor   string         // название для ссылки
	Alias    string         // имя ссылки на исходные данные
	position `yaml:"-"`     // позиция в исходном файле
}

// UnmarshalYAML реализует интерфейс [yaml.Unmarshaler].
func (fs *Fields) UnmarshalYAML(n *yaml.Node) error {
	fs.Anchor = n.Anchor // сохраняем имя ссылки

	// проверяем, что данные не определены через ссылку
	if n.Kind == yaml.AliasNode {
		n = n.Alias         // подставляем оригинальные данные для разбора
		fs.Alias = n.Anchor // запоминаем имя ссылки
	}

	if n.Kind != yaml.MappingNode {
		return NewError(nil, n, "fields must be a YAML mapping: have %v", n.Kind)
	}

	// сохраняем позицию в исходном файле с определением элемента списка
	fs.position = parseSource(n)

	// инициализируем список для хранения описания полей
	count := len(n.Content) / 2
	fs.Fields = make([]Field, 0, count)
	fs.index = make(map[string]int, count)

	// разбираем дерево с описанием полей в формате YAML
	for i := 1; i < len(n.Content); i += 2 {
		var f Field

		nameNode := n.Content[i-1] // информация о названии
		f.Name = nameNode.Value    // сохраняем название запроса
		if f.Name == "" {
			return NewError(nil, nameNode, "field name not defined")
		}

		// проверяем, что такое поле ещё не было определено ранее
		if _, ok := fs.index[f.Name]; ok {
			return NewError(nil, nameNode, "field %q redefined", f.Name)
		}

		// сохраняем позицию в исходном файле с определением элемента списка
		f.position = parseSource(nameNode)

		// разбираем поля с описанием типа
		valueNode := n.Content[i]
		// тип данных поля
		f.Type = valueNode.Value
		if f.Type == "" {
			return NewError(nil, valueNode, "field %q type not defined", f.Name)
		}

		// комментарий
		f.Comment = parseComments(nameNode, valueNode)

		// сохраняем разобранный запрос и его индекс
		fs.Fields = append(fs.Fields, f)
		fs.index[f.Name] = len(fs.Fields) - 1
	}

	return nil
}
