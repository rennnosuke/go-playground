## How to Test

1. Run the following command to start the MySQL container:

```sh
$ init.sh
```

2. Exec test in the following command:

```sh
$ go test ./db/rdb/mysql/conn 
```