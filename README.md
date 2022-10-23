# Webshop Web App

This is an example app demonstrating clean architecture and various tools 

## Run application

To run the app use the provided Makefile.
Use `make run` to start the database and application.

## Testing

There are multiple types of tests:
1. unit tests - can be run anywhere
2. integration tests - needs a docker socket
3. api tests - needs a docker socket + newman cli

To run the unit tests use: `make short`

To run the unit tests + integration tests use: `make test`

To run the api tests use: `make apitest`