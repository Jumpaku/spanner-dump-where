# spanner-dump-where (v0.0.0)


## spanner-dump-where

### Description

Dump data from a Google Cloud Spanner database with specified conditions.
This command allows you to export data from a Spanner database, applying filters and options to control the output.

### Syntax

```shell
spanner-dump-where [<option>]...
```

### Options

* `-bulk-size=<integer>`  (default=`100`):  
  Number of rows to dump in a single batch.  
  This option is used to control the size of the data dump.  

* `-database=<string>`, `-d=<string>`  (default=`""`):  
  Google Cloud Spanner database ID.  
  This option is required.  

* `-from=<string>`  (default=`""`):  
  Table name to dump data from.  
  This option can be specified one or more times.  

* `-instance=<string>`, `-i=<string>`  (default=`""`):  
  Google Cloud Spanner instance ID.  
  This option is required.  

* `-no-data[=<boolean>]`  (default=`false`):  
  If true, do not dump data.  

* `-no-ddl[=<boolean>]`  (default=`false`):  
  If true, do not dump DDL statements.  

* `-project=<string>`, `-p=<string>`  (default=`""`):  
  Google Cloud project ID.  
  This option is required.  

* `-sort[=<boolean>]`  (default=`false`):  
  If true, sort the dump order according to dependency relationships on tables.  
  This option is used to control the order of the dumped data.  

* `-timestamp=<string>`, `-t=<string>`  (default=`""`):  
  Timestamp to use for the dump.  

* `-upsert[=<boolean>]`  (default=`false`):  
  If true, use INSERT OR UPDATE instead of INSERT.  

* `-where=<string>`  (default=`""`):  
  Condition to filter data.  
  This option is required for each -from option.  
  The format is an SQL boolean expression after WHERE clause.  




