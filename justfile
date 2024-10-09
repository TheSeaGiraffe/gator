#set shell := ["zsh", "-cu"]
set dotenv-load

# default task
default:
    just --list --unsorted

# run the app with the specified args
run *args:
    go run . {{args}}

# spin up DB container
db_up:
    docker compose up -d 

# shut down DB container
db_down:
    docker compose down

# connect to the DB using the provided DSN
db_connect:
    psql "$PSQL_DSN"

# perform the specified migration action
migrate action="status":
    goose -dir sql/schema postgres "$PSQL_DSN" {{action}}
