#!/bin/sh

# create .envrc for direnv
cat <<-EOS >> .envrc
	export MYSQL_LOCAL_VERSION=9.0.1
	export MYSQL_LOCAL_HOST=localhost
	export MYSQL_LOCAL_PORT=3306
	export MYSQL_LOCAL_DATABASE=test
	export MYSQL_LOCAL_USER=testuser
	export MYSQL_LOCAL_PASSWORD=testpass
	export MYSQL_LOCAL_ROOT_PASSWORD=my-secret-pw
EOS

# enable direnv
if ! command -v 2>&1 >/dev/null
then
    echo "direnv could not be found"
    exit 1
fi
direnv allow .
