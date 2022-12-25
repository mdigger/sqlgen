package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Error описывает ошибку разбора конфигурации.
type Error struct {
	Message  string     // сообщение об ошибке
	Query    string     // название запроса
	err      error      // оригинальная ошибка
	position `yaml:"-"` // строка и позиция в исходном файле
}

// Error возвращает строку с описанием ошибки.
func (e Error) Error() string {
	message := fmt.Sprintf("[%v] %q %s", e.Source(), e.Query, e.Message)
	if e.err != nil {
		message += ": " + e.err.Error()
	}

	return message
}

// Unwrap возвращает оригинальную ошибку.
func (e Error) Unwrap() error {
	return e.err
}

// NewError формирует и возвращает описание ошибки при разборе запроса.
func NewError(err error, n *yaml.Node, format string, args ...any) error {
	qerr := Error{
		Message: fmt.Sprintf(format, args...),
		err:     err,
	}

	if n != nil {
		qerr.Query = n.Value
		qerr.position = parseSource(n)
	}

	return qerr
}
