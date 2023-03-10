# sqlgen

Генератор библиотеки с SQL-запросами для Golang.

## Описание запроса

Вы описываете SQL запросы и параметры, которые в них используются, а sqlgen генерирует библиотеку Golang для работы с ними.

#### Пример:

```yaml
# Return user information.
get user: 
  type: one # only one row
  sql: |-
    select name, age, email
    from users
    where id = ?
  in:
    id: string # user id
  out: # User info.
    name: string # user name
    age: uint # age
    email: sql.NullString # email
```

Сначала идёт название запроса (в нашем примере это `get user`).
В одном файле может быть описано сразу несколько запросов.

- **`type`** -- тип запроса:
  - `many` -- позволяет обработать несколько возвращаемых записей 
  - `one` -- возвращает только одну запись из базы данных
  - `exec` -- не возвращает значений (кроме возможной ошибки выполнения)
  - `id` -- возвращает идентификатор записи (`int64`), сгенерированный сервером
  - `affected` -- возвращает количество обработанных запросом записей
  - `exist` -- возвращает ошибку `sql.ErrNoRows`, если ни одна запись не была обработана (`affected == 0`)
- **`sql`** -- описание запроса в формате SQL
- **`in`** -- список параметров запроса с их названиями и типами
  - `name`: `type`
  - ...
- **`out`** -- так же содержит список параметров, но уже с описанием возвращаемых значений.
  - `name`: `type`
  - ...

Поддержка разбора [нескольких одновременных запросов](https://pkg.go.dev/database/sql#Rows.NextResultSet) не реализована и пока не планируется.

Комментарии из описания, по-возможности, переносятся в сгенерированный код, поэтому ими не стоит пренебрегать.
Многострочный комментарий можно писать перед параметром, группы параметров или названием запроса. Короткие комментарии -- сразу после определения в той же строке. Последние хорошо подходят в качестве комментариев для описания параметров.


### Поля запроса

Порядок описания параметров должен соответствовать тому, как они описаны в запросе.

В качестве типов параметров поддерживаются стандартные типы golang `string`, `int`, `uint`, `bool`, `float32` и так далее. 

Кроме этого, сразу добавлена возможность использования так же:
- [`sql.NullString`](https://pkg.go.dev/database/sql#NullString), 
- [`sql.NullTime`](https://pkg.go.dev/database/sql#NullTime), 
- [`sql.NullBool`](https://pkg.go.dev/database/sql#NullBool), 
- [`sql.NullByte`](https://pkg.go.dev/database/sql#NullByte), 
- [`sql.NullInt16`](https://pkg.go.dev/database/sql#NullInt16),
- [`sql.NullInt32`](https://pkg.go.dev/database/sql#NullInt32),
- [`sql.NullInt64`](https://pkg.go.dev/database/sql#NullInt64),
- [`sql.NullFloat64`](https://pkg.go.dev/database/sql#NullFloat64),
- [`sql.RawBytes`](https://pkg.go.dev/database/sql#RawBytes),
- [`json.RawMessage`](https://pkg.go.dev/encoding/json#RawMessage),
- [`time.Time`](https://pkg.go.dev/encoding/time#Time),

Поддержка [именованных параметров](https://pkg.go.dev/database/sql#NamedArg) пока не планируется, потому что они не поддерживаются в MySQL и потребуют некоторой дополнительной логики для генератора. Аналогично, и [sql.Out](https://pkg.go.dev/database/sql#Out).


### Сторонние библиотеки с типами данных

Если вы хотите использовать другие типы, поддерживающие интерфейс `sql.Scanner` и `drive.Valuer`, то необходимо явно указать при использовании генератора на использование этих библиотек:

```shell
$ sqlgen generate --import github.com/gofrs/uuid
```
В том случае, если название библиотеки не совпадает с используемым префиксом данных или  использует суффикс с версией, необходимо так же указать и то, какой префикс данных она обслуживает через двоеточие перед названием:

```shell
$ sqlgen generate --import uuid:github.com/vgarvardt/pgx-google-uuid/v5
```

### Повторяющиеся списки параметров

Иногда список возвращаемых полей повторяется в нескольких запросах. Чтобы не дублировать его вручную, можно использовать синонимы при описании.

```yaml
get user:
  type: one
  sql: |-
    select id, name, age, email
    from users
    where id = ?
  in:
    id: string
  out: &user # <- named param list
    id: string
    name: string # user name
    age: uint # age
    email: sql.NullString # email
```

В примере выше мы задали имя для списка исходящих параметров (`&user`).
В дальнейшем это имя используется вместо повторения списка параметров (`*user`):

```yaml
get all users:
  type: many
  sql: |-
    select id, name, age, email
    from users
  out: *user # <- using named param list
```


## Генерация 

Данная команда генерирует код библиотеки с SQL-запросами. По умолчанию сгенерированные файлы записываются в текущий каталог. С помощью флага out можно явно указать каталог для генерации файлов:

```shell
$ sqlgen generate --out ./database
```

В качестве названия библиотеки по умолчанию используется название каталога. Если вы хотите это изменить, то задайте новое имя библиотеки через флаг name:

```shell
$ sqlgen generate --out ./database --name db
```

Изначально генератор автоматически поддерживает стандартные типы данных, определенные в языке golang: string, int, int64 и так далее. Кроме этого по умолчанию поддерживаются следующие пакеты:

- `database/sql`
- `encoding/json`
- `time`

Если в описании типов данных вы хотите использовать сторонние библиотеки, то это необходимо явно указать, задав их с помощью флага import:

```shell
$ sqlgen generate --import github.com/gofrs/uuid
```

Префикс библиотеки определяется по последнему элементу в пути к пакету и не содержит в себе другого кода для определения реального названия. Поэтому, если используемый вами префикс отличается, то его необходимо явно указать перед названием пакета через двоеточие:

```shell
$ sqlgen generate --import uuid:github.com/vgarvardt/pgx-google-uuid/v5
```

