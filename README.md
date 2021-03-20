## Raccoon Stats API
>> Stats tracking app

Made with Go, GraphQL ([gqlgen](https://github.com/99designs/gqlgen)) and PostgreSQL (10.15)

### Run

Create database schema
```
psql -h 12.13.14.15 -p 5432 -d DB_NAME -U DB_USER -a -f scripts/create_database.sql -W
```

Rename **.env-example** to **.env** and set values

Run with
`make run`

### Docker
Rename **.env-example** to **.env** and set values

Build
`docker build . -t rs_api_docker:latest`

Run
`docker run -p 8080:4200 rs_api_docker`

Logs
`docker ps -a docker logs ce62844970f8`

## Development
1. Edit GraphQL-schema in `graph/schema/schema.graphqls`

2. Run `make generate` which is alias to `go run github.com/99designs/gqlgen generate`

3. Copy new resolver prototypes from generated file `graph/schema.resolvers.go` to more suitable place and implement resolvers
