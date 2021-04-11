## Raccoon Stats API
>> Stats tracking app

Made with Go, GraphQL ([gqlgen](https://github.com/99designs/gqlgen)) and PostgreSQL (10.15)

### Run

Rename **.env-example** to **.env** and set values

Create database schema
```
psql -h 12.13.14.15 -p 5432 -d DB_NAME -U DB_USER -a -f scripts/create_database.sql -W
```

Run with
`make run`

### Docker
Rename **.env-example** to **.env** and set values

Compose
`docker-compose up`

## Development
1. Edit GraphQL-schema in `graph/schema/schema.graphqls`

2. Run `make generate` which is alias to `go run github.com/99designs/gqlgen generate`

3. Copy new resolver prototypes from generated file `graph/schema.resolvers.go` to more suitable place and implement resolvers
