set dotenv-path := "./docker/.env"
project_path := "./orgstruct"
project_prefix := "orgstruct_"
migrate_volume_name := "migrations_data"


#################################
### app

# run main.go by default with the env production parameter
run:
    cd {{project_path}} && go run ./cmd/main.go -env dev

# run main.go with the env parameter
run-env env:
    cd {{project_path}} && go run ./cmd/main.go -env {{env}}

# rename go mod file: github.com/<user-name>/<name>
go-mod-edit name:
    go mod edit -module {{ name }} {{project_path}}/go.mod

#################################
### docker 

# run docker postgres db, build image orgstruct-goose-migrate, implement migrations, remove migrations container
docker-deploy: docker-up-db docker-build-migrate migrate-init migrate-up

# start db container
docker-start-db:
    docker compose -f ./docker/docker-compose.yml start orgstruct-postgres

# create containers applies migrations on database and remove migration container
docker-up-db:
    docker compose -f ./docker/docker-compose.yml up -d orgstruct-postgres

#################################
### migrate

# create or update image with an argument indicating the path to migration from justfile
docker-build-migrate path="{{project_path}}/migrations":
    docker build -t orgstruct-goose-migrate:1.0.0 \
    --build-arg MIGRATIONS_PATH={{path}} \
    -f ./docker/Dockerfile.migrate .

# create migration file with migrations_name
migrate-create migration_name:
    cd {{project_path}} \
    && go run github.com/pressly/goose/v3/cmd/goose@latest \
    -dir ./migrations create {{migration_name}} \
    sql

# add or update migrations in volume for load in migrate
migrate-init: 
    docker compose -f ./docker/docker-compose.yml run --rm \
        -v "{{project_path}}/migrations:/source:ro" \
        -v {{migrate_volume_name}}:/migrations \
        migrate-init \
        cp -r /source/. /migrations/

# use: up, up-to, down, down-to, status, redo, reset. Launch migrate container, implement migrations and remove container. 
migrate cmd  *args:
    docker compose -f ./docker/docker-compose.yml run --rm \
    -v {{project_prefix}}{{migrate_volume_name}}:/migrations \
    orgstruct-migrate -dir=/migrations {{ cmd }} {{args}}

# migrate up
migrate-up:
    migrate up