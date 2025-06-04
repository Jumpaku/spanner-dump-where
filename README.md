spanner-dump-where
===

spanner-dump-where is a command line tool for conditionally exporting a Google Cloud Spanner database in SQL format.
Exported databases can be imported to Google Cloud Spanner with [spanner-cli](https://github.com/cloudspannerecosystem/spanner-cli).

```sh
# Export
$ spanner-dump-where -project=${PROJECT} -instance=${INSTANCE} -database=${DATABASE} \
  -sort \
  -bulk-size=10 \
  -from=User -where='Age > 20' \
  -from=UserItem -where='UserId = "1"' \
  > data.sql

# Import
$ spanner-cli -p ${PROJECT} -i ${INSTANCE} -d ${DATABASE} < data.sql
```

spanner-dump-where enhances spanner-dump with the following features:
- It can specify conditions to filter data using an SQL boolean expression.
- It can sort the dump order according to dependency relationships, such as interleave and foreign keys.
- It can use INSERT OR UPDATE instead of INSERT.

spanner-dump-where is a fork of https://github.com/cloudspannerecosystem/spanner-dump .

## Limitations

- This tool does not ensure consistency between the database schema (DDL) and data. Therefore, you should avoid making changes to the schema while running this tool.

## Install

```
go install github.com/Jumpaku/spanner-dump-where@latest
```

## Usage

```
    spanner-dump-where

    Description:
        Dump data from a Google Cloud Spanner database with specified conditions.
        This command allows you to export data from a Spanner database, applying filters and options to control the output.

    Syntax:
        $ spanner-dump-where [<option>]...

    Options:
        -bulk-size=<integer>  (default=100):
            Number of rows to dump in a single batch.
            This option is used to control the size of the data dump.

        -database=<string>, -d=<string>  (default=""):
            Google Cloud Spanner database ID.
            This option is required.

        -from=<string>  (default=""):
            Table name to dump data from.
            This option can be specified one or more times.

        -instance=<string>, -i=<string>  (default=""):
            Google Cloud Spanner instance ID.
            This option is required.

        -no-data[=<boolean>]  (default=false):
            If true, do not dump data.

        -no-ddl[=<boolean>]  (default=false):
            If true, do not dump DDL statements.

        -project=<string>, -p=<string>  (default=""):
            Google Cloud project ID.
            This option is required.

        -sort[=<boolean>]  (default=false):
            If true, sort the dump order according to dependency relationships on tables.
            This option is used to control the order of the dumped data.

        -timestamp=<string>, -t=<string>  (default=""):
            Timestamp to use for the dump.

        -upsert[=<boolean>]  (default=false):
            If true, use INSERT OR UPDATE instead of INSERT.

        -where=<string>  (default=""):
            Condition to filter data.
            This option is required for each -from option.
            The format is an SQL boolean expression after WHERE clause.
```
