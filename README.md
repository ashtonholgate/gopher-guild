# Gopher Guild

Gopher Guild is an example monorepo based microservice architecture which currently contains a single example microservice called `gopherService`. `gopherService` is a Go-based web application that provides an API for managing gopher entities which uses the Gin web framework. Separately to this a Postgres DB is used for persisting the information created by the `gopherService`.

## Purpose

The purpose of this project was in order to learn idiomatic Go, and to that end the project has been architected in the following ways:

- Extensive use of dependency injection so as to make the code easily testable.
- Having services define their own dependencies, so that minimally slim mocks can be given when testing.
- Services have minimal number of jobs and responsibilities to make the project as modular as possible whilst also adhering to the practice of separation of concerns.
- The main.go file has been slimmed down as much as possible, and all logic has been put in sub-packages instead, so as to maximise testability
- Use of custom errors at the point of error, which are wrapped with `fmt.Errorf` methods at each layer that the error passes through to add context, to then be unwrapped at the controller level with an error type switch so that the full context can be logged for development purposes, whereas a nice looking error can be returned to the client that made the request.

## Architecture

The gopherService is composed out of services, each of which serve distinct purposes:

- Controller: for analysing requests, making sure that they conform to requirements, and mapping it before passing on to a command or query service where applicable.
- Command / Query Service: for orchestrating what needs to happen in order to comply with requests made by the controller. This will typically include making multiple requests in a specific order.
- Repository Service: for interacting directly with the database, and returning suitably mapped data back to the command / query services.

## Setup and Installation

This project is intended to be run using Docker Compose. As such please ensure that you have an instance of Docker running on your machine, and have docker compose available.

Once installed and set up, run the command `docker compose up` in your terminal, scoped to this directory, in order to create both a running instance of the gopherService in addition to a database to use.

The project will update live if you change the code, leading to easy development.

If it is your first time running `docker compose up`, or if you have previously destroyed your instance of the database either with a `docker compose down -v` command, or by other means, you must run the migrations. I use [Go Migrate](https://github.com/golang-migrate/migrate) to do so. Once installed, you can run the command `migrate -path ./migrations -database "postgres://user:password@localhost:5432/db?sslmode=disable" up` to run all of the migrations found in the `./migrations` directory.

## API Specification

You can find an OpenAPI specification for the gopher service in `./gopherService/apiSpecification.yaml`

## Testing

When scoped to the gopherService, you can run tests for all packages by running `go test ./...` in the terminal.

## Environment Variables

This project, as it is an example project, has the environment variables committed to the repo.
