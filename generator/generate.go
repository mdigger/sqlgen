package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mdigger/sqlgen/config"
)

var (
	//go:embed generator.tmpl
	queryTemplates string
	// tmpl содержит разобранные шаблоны для генерации кода.
	tmpl = template.Must(template.New("").Funcs(funcMap).Parse(queryTemplates))
)

// Generator описывает данные генератора.
type Generator struct {
	Name    string // название
	Version string // версия
	Package string // название пакета

	imports map[string]string // список поддерживаемых импортов пакетов по префиксам
}

// New возвращает новый генератор с заданным именем библиотеки для генерации.
// Опционально указываются дополнительный пакеты для импорта.
func New(name string, packages ...string) Generator {
	if name == "" || name == "." {
		name, _ = os.Getwd() // используем текущий каталог
	}

	name = path.Base(filepath.ToSlash(name)) // отделяем имя от пути
	switch name {
	case "", ".", "/":
		name = "database"
	}

	// используемые пакеты по умолчанию
	imports := map[string]string{
		"sql":  "database/sql",
		"time": "time",
		"json": "encoding/json",
	}

	// добавляем дополнительные пакеты
	for _, item := range packages {
		if item == "" {
			continue
		}

		if idx := strings.IndexByte(item, ':'); idx > 0 {
			imports[item[:idx]] = item[idx+1:] // задан префикс
		} else {
			// TODO: добавить разбор версии пакета
			imports[path.Base(item)] = item
		}
	}

	return Generator{
		Name:    "github.com/mdigger/sqlgen",
		Version: "v0.1.0",
		Package: name,
		imports: imports,
	}
}

// Query генерирует и возвращает код для работы с запросами.
func (g Generator) Query(source string, queries []config.Query) ([]byte, error) {
	// определяем список библиотек, используемых в запросах, для импорта
	imports, err := g.getImports(queries)
	if err != nil {
		return nil, err
	}

	// формируем данные для использования в шаблоне
	data := struct {
		Generator                   // информация о генераторе
		Source    string            // название и путь исходного файла с данными
		Imports   map[string]string // список импортируемых библиотек
		Queries   []config.Query    // список запросов
	}{
		Generator: g,
		Source:    source,
		Imports:   imports,
		Queries:   queries,
	}

	// генерируем и возвращаем код для обработки запроса
	return generate("generate queries", data)
}

// DB возвращает сгенерированный код с описанием библиотеки запросов.
func (g Generator) DB() ([]byte, error) {
	// формируем данные для использования в шаблоне
	data := struct {
		Generator        // информация о генераторе
		Source    string // для основного модуля исходный файл не задаётся
	}{Generator: g}

	// генерируем и возвращаем код с основным описанием библиотеки
	return generate("generate db", data)
}

// generate генерирует код с использованием шаблона name и параметров data.
// Возвращает форматированный сгенерированный код.
func generate(name string, data any) ([]byte, error) {
	// генерируем код основного файла на основании шаблона
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	// форматируем код согласно принятым правилам
	source := buf.Bytes()
	formatted, err := format.Source(source)
	if err != nil {
		return source, fmt.Errorf("format: %w", err)
	}

	return formatted, nil
}

// getImports возвращает список библиотек для импорта.
func (g Generator) getImports(qs []config.Query) (map[string]string, error) {
	// определяем, какие библиотеки нужно импортировать
	used := make(map[string]string, len(g.imports))

	// вспомогательная функция для формирования списка используемых модулей
	getImports := func(typeName, queryName string) error {
		if idx := strings.IndexByte(typeName, '.'); idx > 0 {
			prefix := typeName[:idx]

			if lib, ok := g.imports[prefix]; ok {
				if strings.HasSuffix(lib, prefix) {
					prefix = "" // не используем префикс, если он явно указан в конце пакета
				}

				used[lib] = prefix
			} else {
				return fmt.Errorf("unknown param package prefix %q in query %q",
					prefix, queryName)
			}
		}

		return nil
	}

	// проходим по всем параметрам (входящим и исходящим) всех запросов и
	// выбираем используемые библиотеки
	for _, q := range qs {
		for _, t := range q.Params.Fields {
			if err := getImports(t.Type, q.Name); err != nil {
				return nil, err
			}
		}

		for _, t := range q.Out.Fields {
			if err := getImports(t.Type, q.Name); err != nil {
				return nil, err
			}
		}
	}

	return used, nil
}
