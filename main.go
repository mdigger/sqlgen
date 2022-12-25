package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mdigger/sqlgen/config"
	"github.com/mdigger/sqlgen/generator"
	"github.com/mdigger/wordwrap"
	"github.com/urfave/cli/v3"
)

func init() {
	log.SetFlags(0) // убираем все флаги вывода в лог
}

var (
	module  = "github.com/mdigger/sqlgen"
	version = "v0.1.0"
)

func main() {
	// конфигурация приложения, поддерживаемых команд и флагов
	app := &cli.App{
		Name:           "sqlgen",
		Usage:          "Golang SQL queries library generator.",
		Version:        version,
		Description:    helpString(appDescription),
		DefaultCommand: "generate",
		Commands: []*cli.Command{{
			Name:        "generate",
			Usage:       "Generate Golang library",
			Description: helpString(generateDescription),
			Action:      generateCmd,
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:    "out",
					Usage:   "`path` to writing",
					Aliases: []string{"o"},
				},
				&cli.StringFlag{
					Name:    "name",
					Usage:   "set `library` name",
					Aliases: []string{"n"},
				},
				&cli.StringSliceFlag{
					Name:    "import",
					Usage:   "import `package`",
					Aliases: []string{"i"},
				},
			},
			// }, {
			// 	Name:  "format",
			// 	Usage: "Format source query files",
		}},
		Authors: []*cli.Author{{
			Name:  "Dmitry Sedykh",
			Email: "sedykh@gmail.com",
		}},
		EnableBashCompletion: true,
		Suggest:              true,
	}

	// запускаем приложение
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}

// generateCmd выполняет команду генерации кода библиотеки.
func generateCmd(c *cli.Context) error {
	// список аргументов для выполнения команды
	args := c.Args().Slice()
	if len(args) == 0 {
		args = []string{"*.yaml"}
	}

	// формируем список файлов с описанием запросов
	files := make(map[string]struct{})
	for _, arg := range args {
		// добавляем маску для выбора файлов, если не указан конкретный файл
		if filepath.Ext(arg) == "" {
			arg = filepath.Join(arg, "*.yaml")
		}

		// получаем список имен файлов для обработки
		matches, err := filepath.Glob(arg)
		if err != nil {
			return fmt.Errorf("match files %q: %w", arg, err)
		}

		// добавляем в список файлов на обработку
		for _, file := range matches {
			files[file] = struct{}{}
		}
	}

	// проверяем, что есть файлы с описанием запросов
	if len(files) == 0 {
		return errors.New("the files with the description of the request were not found")
	}

	outFolder := c.Path("out") // каталог для записи файлов
	// создаём каталог для сохранения сгенерированных файлов, если его нет
	if outFolder != "" && outFolder != "." {
		if _, err := os.Stat(outFolder); os.IsNotExist(err) {
			log.Println("creating output folder:", outFolder)
			if err := os.MkdirAll(outFolder, 0o750); err != nil {
				return fmt.Errorf("output folder %q: %w", outFolder, err)
			}
		}
	}

	name := c.String("name") // название пакета
	if name == "" {
		name = outFolder
	}

	// инициализируем генератор кода с заданным именем и списком импортируемых библиотек
	generator := generator.New(name, c.StringSlice("import")...)

	// обрабатываем все файлы из нашего списка
	for file := range files {
		// разбираем описание запроса из файла
		qs, err := config.Parse(file)
		log.Println("parsing file:", file)
		if err != nil {
			return fmt.Errorf("parse %q: %w", file, err)
		}

		// получаем сгенерированный код с описанием запросов
		data, err := generator.Query(file, qs.Queries)
		log.Println("generating code:", len(data), "bytes")
		if err != nil {
			os.Stderr.Write(data)
			return fmt.Errorf("generator %q: %w", file, err)
		}

		// формируем новое имя файла и записываем в него получившийся код
		destination := filepath.Join(outFolder, filepath.Base(file))
		destination = destination[:len(destination)-len(filepath.Ext(destination))]
		destination += ".sql.go"
		log.Println("saving file:", destination)
		if err = os.WriteFile(destination, data, 0o666); err != nil {
			return fmt.Errorf("save file %q: %w", destination, err)
		}
	}

	// генерируем код инициализации библиотеки
	data, err := generator.DB()
	log.Println("generating package code:", len(data), "bytes")
	if err != nil {
		return fmt.Errorf("generate main: %w", err)
	}

	// записываем код библиотеки в файл
	destination := filepath.Join(outFolder, "db.go")
	log.Println("saving package code:", destination)
	if err = os.WriteFile(destination, data, 0o666); err != nil {
		return fmt.Errorf("save main %q: %w", destination, err)
	}

	log.Println("generation completed!")
	return nil
}

// helpString возвращает текст с переносом по строкам.
func helpString(s string) string {
	return wordwrap.String(s, 72)
}

const (
	appDescription      = `The description of SQL queries and the parameters used in them is done using YAML files. Based on these descriptions, sqlgen generates a library to work with these queries.`
	generateDescription = `This command generates the golang library with SQL queries.
	
By default, the generated files are written to the current directory. Using the flag "out" you can explicitly specify a directory for generating files:
	sqlgen generate --out ./database

The default library name is the directory name. If you want to change this, then set a new library name via the "name" flag:
	sqlgen generate --out ./database --name db

Initially, the generator automatically supports standard data types defined in the golang language: string, int, int64, and so on. In addition, the following packages are supported by default:
	database/sql, encoding/json, time

If you want to use third-party libraries in the description of data types, then this must be explicitly specified by setting them using the "import" flag:
	sqlgen generate --import github.com/gofrs/uuid

The library prefix is determined by the last element in the package path and does not contain any other code to determine the actual name. Therefore, if the prefix you are using is different, then it is necessary explicitly specify a colon before the package name:
	sqlgen generate --import uuid:github.com/vgarvardt/pgx-google-uuid/v5`
)
