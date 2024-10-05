## How to Test

1. Run the following command to start the MySQL container:

```bash
$ docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:9.0.1
$ docker exec -it <container_id> /bin/bash
```