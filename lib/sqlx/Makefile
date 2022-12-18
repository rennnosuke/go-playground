.PHONY: build, run
build:
	docker build -t docker-mysql-5-7:1 .
run:
	docker run --env-file ".env" --name docker-mysql-5-7 -d -v $$PWD/db-5-7:/var/lib/mysql -p 23306:3306 docker-mysql-5-7:1