BINARY_NAME = rs_api

build:
	go build -o bin/$(BINARY_NAME) cmd/rs/*.go

run: build
	./bin/$(BINARY_NAME)

# Generate data fromGraphql-schema
generate:
	@go run github.com/99designs/gqlgen generate
	@rm graph/resolver.go
	# Copy prototype form this file!
	# rm graph/schema.resolvers.go

# Init database
database:
	@psql -h $(RS_DB_HOST) -p $(RS_DB_PORT) -d $(RS_DB_NAME) -U $(RS_DB_USER) -a -f scripts/create_database.sql -W

# Insert sample data
seed:
	@psql -h $(RS_DB_HOST) -p $(RS_DB_PORT) -d $(RS_DB_NAME) -U $(RS_DB_USER) -a -f scripts/seed_database.sql -W

# Deploy to server
deploy: build
	cp bin/$(BINARY_NAME) ./ansible/$(BINARY_NAME)
	cp .env ./ansible/.env
	(cd ansible && ansible-playbook deploy.yml)
	rm ./ansible/.env
	rm ./ansible/$(BINARY_NAME)

test:
	go test -v ./...

# First run ?
#init-gql:
#	go run github.com/99designs/gqlgen init
