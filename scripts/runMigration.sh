#!/bin/bash

MIGRATE_PATH="${GOPATH:-$HOME/go}/bin/migrate"

check_migrate() {
    if ! command -v $MIGRATE_PATH &> /dev/null; then
        echo "migrate tool not found. Installing..."
        go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    fi
}

source .env

echo "DB_USER: $DB_USER"
echo "DB_PASSWORD: [hidden]" 
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "DB_NAME: $DB_NAME"

check_migrate

if [ -z "$DB_PASSWORD" ]; then
    DB_URL="mysql://$DB_USER@tcp($DB_HOST:$DB_PORT)/$DB_NAME"
else
    DB_URL="mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME"
fi
echo "Constructed DB_URL: $DB_URL"

$MIGRATE_PATH -path database/migrations -database "$DB_URL" up
