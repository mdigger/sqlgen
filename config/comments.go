package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Comment описывает комментарии.
type Comment []string

// IsMultiline возвращает true, если комментарий многострочный.
func (c Comment) IsMultiline() bool {
	return len(c) > 1
}

// Format возвращает комментарий, добавляя префикс к каждой строке.
// Возвращает пустую строку, если комментарии не заданы.
func (c Comment) Format(prefix string) string {
	if len(c) == 0 {
		return ""
	}

	// каждая строка комментария начинается с префикса
	return prefix + strings.Join(c, "\n"+prefix)
}

// GoString возвращает строку с представлением комментариев для языка Golang.
func (c Comment) GoString() string {
	return c.Format("// ")
}

// parseComments используется для получения строки комментария из описания YAML-нод.
// Можно указать сразу список нод, которые будут перебираться в указанном порядке, пока не найдётся первый
// не пустой комментарий.
//
// Если комментарий состоит из нескольких строк, то каждая строка представлена отдельным элементом массива.
// Автоматически удаляет префикс # из текста комментариев.
func parseComments(n ...*yaml.Node) Comment {
	for _, n := range n {
		comment := n.HeadComment
		if comment == "" {
			comment = n.LineComment
		}

		if comment == "" {
			continue
		}

		comments := strings.Split(comment, "\n")
		for i, comment := range comments {
			if comment != "" && comment[0] == '#' {
				comments[i] = strings.TrimSpace(comment[1:])
			}
		}

		return comments
	}

	return nil
}
