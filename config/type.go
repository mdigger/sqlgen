package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Type описывает тип запроса.
type Type uint8

// Предопределенный список поддерживаемых типов запросов.
// В зависимости от типа, используются разные методы запроса и возвращаемые данные.
const (
	TypeMany     Type = iota // список записей
	TypeOne                  // только одну запись
	TypeExec                 // выполнение запроса без возврата ответа
	TypeAffected             // выполнение и возврат количества затронутых изменениями записей
	TypeExist                // возвращает true, если хотя бы одна запись изменена
	TypeRowID                // возвращает сгенерированный сервером идентификатор записи
)

// String возвращает строку, описывающую идентификатор типа запроса.
func (qt Type) String() string {
	switch qt {
	case TypeMany:
		return "many"
	case TypeOne:
		return "one"
	case TypeExec:
		return "exec"
	case TypeAffected:
		return "affected"
	case TypeExist:
		return "exist"
	case TypeRowID:
		return "id"
	default:
		return ""
	}
}

// MarshalYAML поддерживает интерфейс [yaml.Marshaler].
func (qt Type) MarshalYAML() (any, error) {
	return qt.String(), nil
}

// parseType разбирает строковое представление типа.
func parseType(s string) (Type, error) {
	switch strings.ToLower(s) {
	case "many":
		return TypeMany, nil
	case "one":
		return TypeOne, nil
	case "exec":
		return TypeExec, nil
	case "affected":
		return TypeAffected, nil
	case "id", "row_id", "rowid":
		return TypeRowID, nil
	case "exist":
		return TypeExist, nil
	default:
		return TypeMany, fmt.Errorf("unsupported query type: %v", s)
	}
}

// UnmarshalYAML поддерживает интерфейс [yaml.Unmarshaler].
func (qt *Type) UnmarshalYAML(n *yaml.Node) error {
	var err error
	*qt, err = parseType(n.Value)
	if err != nil {
		return NewError(err, n, err.Error()) // оборачиваем ошибку
	}

	return nil
}
