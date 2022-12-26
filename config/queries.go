package config

import (
	"gopkg.in/yaml.v3"
)

// Queries содержит список описаний запросов.
type Queries struct {
	Queries []Query        // список запросов
	index   map[string]int // индекс запросов по их заголовку
}

// Count возвращает количество определенных запросов.
func (qs Queries) Count() int {
	return len(qs.Queries)
}

// Names возвращает список имён запросов.
func (qs Queries) Names() []string {
	if len(qs.Queries) == 0 {
		return nil
	}

	list := make([]string, len(qs.Queries))
	for i := range qs.Queries {
		list[i] = qs.Queries[i].Name
	}

	return list
}

// Get возвращает ссылку на описание запроса с указанным именем.
func (qs Queries) Get(name string) *Query {
	if idx, ok := qs.index[name]; ok {
		return &qs.Queries[idx]
	}

	return nil
}

// UnmarshalYAML реализует интерфейс [yaml.Unmarshaler].
func (qs *Queries) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return NewError(nil, n, "query must be a YAML mapping: have %v", n.Kind)
	}

	// инициализируем список запросов и индекс
	count := len(n.Content) / 2
	qs.Queries = make([]Query, 0, count)
	qs.index = make(map[string]int, count)

	// разбираем дерево с описанием запросов в формате YAML
	for i := 1; i < len(n.Content); i += 2 {
		var q Query

		nameNode := n.Content[i-1] // информация о названии
		q.Name = nameNode.Value    // сохраняем название запроса

		// проверяем, что имя запроса уникально и ещё не использовалось
		if _, ok := qs.index[q.Name]; ok {
			return NewError(nil, n, "query redefined")
		}

		// комментарий
		q.Comment = parseComments(nameNode)

		// заполняем информацию о запросе
		if err := n.Content[i].Decode(&q); err != nil {
			return NewError(err, n, "query decode")
		}

		// сохраняем разобранный запрос и его индекс
		qs.Queries = append(qs.Queries, q)
		idx := len(qs.Queries) - 1
		qs.index[q.Name] = idx
	}

	return nil
}

// // MarshalYAML поддерживает интерфейс [yaml.Marshaler].
// func (qs Queries) MarshalYAML() (any, error) {
// 	return nil, nil
// }
