# db-schema-doc

A CLI tool to generate database definition doc for doc document from a database table.

## Usage

```sh
db-schema-doc -D my_database -t books
```

And then, you'll get in stdout:

```
| FIELD |    TYPE     | NULL | DEFAULT |       COMMENT        |
|-------|-------------|------|---------|----------------------|
| id    | int(11)     | NO   |         |                      |
| title | varchar(32) | NO   |         | useful title comment |
| price | int(11)     | YES  |         | useful price comment |
```

### Details

```
"db-schema-doc" is a CLI tool to generate a table definition doc
from an existing database table.

Positional arguments specify columns to show.

Usage:
  db-schema-doc [flags]

Flags:
      --config string     config file (default $HOME/.db-schema-doc.yml)
  -D, --database string   database name
  -f, --full              shows all columns if true
  -h, --help              help for db-schema-doc
  -p, --password string   password for database
  -t, --table string      table name
  -u, --username string   username for database (default "root")
```

## TODO

- Convert types to readable forms: `int(11)` -> `Integer`, `varchar(255)` -> `string`.
- Default database configured in YAML.
- Support format of Confluence wiki markup.
- Take columns as flags, not positional arguments.
- Add tests.
