## How to Test

1. Run the following command to start the MySQL container:

```bash
$ docker run --name some-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=my-secret-pw \
  -e MYSQL_DATABASE=test \
  -e MYSQL_HOST=localhost \
  -e MYSQL_PORT=3306 \
  -e MYSQL_USER=testuser \
  -e MYSQL_PASSWORD=testpass \
  -d mysql:9.0.1
```

2. Exec test in the following command:

```bash
$ go test ./db/rdb/mysql/conn 
```