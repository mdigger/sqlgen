package generator

import (
	"go/token"
	"strings"
	"text/template"
	"unicode"
)

// funcMap регистрирует функции для использования в шаблонах.
var funcMap = template.FuncMap{
	"name":   publicName,     // конвертирует строку в название экспортируемого типа
	"param":  param,          // проверяет название параметра
	"escape": escapeBacktick, // экранирует символ "`"
}

// name конвертирует название запроса в название функции golang.
func publicName(s string) string {
	s = name(s, true)

	// подменяем некоторые используемые нами названия параметров
	if s == "Queries" {
		return s + "_"
	}

	return s
}

// param возвращает название параметра.
func param(s string) string {
	// подменяем некоторые используемые нами названия параметров
	switch s {
	case "ctx", "f", "q", "row", "rows", "err", "result", "out", "Queries":
		return s + "_"
	default:
		return name(s, false)
	}
}

// name приводит строку к формату названия в golang.
// Параметр public влияет на заглавную первую букву в имени.
func name(s string, public bool) string {
	// проверяем использование ключевых слов
	if !public && token.IsKeyword(s) {
		return "_" + s
	}

	var (
		buf            strings.Builder
		first          = true // первый символ в названии
		capitalizeNext = public
	)

	for _, r := range s {
		switch {
		case first:
			first = false // больше первого символа при формировании не случится
			// чтобы гарантированно начиналось с буквы, добавляем подчёркивание, если это не так
			if !unicode.IsLetter(r) && r != '_' {
				buf.WriteRune('_')
			} else if public {
				r = unicode.ToTitle(r) // приводим к верхнему регистру
				capitalizeNext = false // следующую букву приводить не требуется
			}

		case !unicode.In(r, unicode.Letter, unicode.Number):
			capitalizeNext = true // символ не является ни буквой, ни цифрой -- сделать следующий заглавной буквой

			continue // пропускаем текущий символ

		case capitalizeNext:
			r = unicode.ToTitle(r) // приводим символ к верхнему регистру
			capitalizeNext = false // сбрасываем флаг, что следующий символ должен быть заглавным
		}

		buf.WriteRune(r) // сохраняем текущий символ
	}

	name := buf.String()
	if strings.HasSuffix(name, "Id") {
		name = name[:len(name)-2] + "ID" // скорее всего это идентификатор
	}

	return name
}

func escapeBacktick(s string) string {
	return strings.Replace(s, "`", "`+\"`\"+`", -1)
}
