# TODO

- [x] использовать комментарии из описания при генерации кода
- [ ] генерировать нормальные ошибки с детальной информацией о проблемном месте в YAML
- [x] выводить в консоль этапы генерации и информацию об обрабатываемых файлах
- [x] поддержка синонимов и ссылок при описании списков полей запросов
- [x] не дублировать код с описанием структуры при использовании синонимов
- [ ] добавить автоматическое форматирование файла с описанием запросов
- [ ] выводить разницу (diff) в случае возможности изменения форматирования запроса
- [ ] разбирать SQL запрос и делать на базе этого дополнительные проверки:
  - [ ] количество описанных входящих параметров должно соответствовать количеству в запросе
  - [ ] предупреждать, если используется `SELECT *`, что это небезопасный способ возврата данных в случае изменения таблицы с данными
  - [ ] определять тип запроса (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) и проверять, что он соответствует типу, указанному в запросе; ругаться на другие типы запросов, что они не поддерживаются
- [ ] проверка корректности описания типов входящих и исходящих параметров
- [ ] рассмотреть возможность поддержки запросов с параметрами в SQL `IN (?)`.