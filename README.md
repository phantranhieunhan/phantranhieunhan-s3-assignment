# Friend Management System
This project provide some features to manage your friends such as
- Connect the friends, list your friends, list your mutual friends to somebody
- Subscribe the friends to receive any updates from them, otherwise you can block to deny any updates or connection from someone you need

## Summary
- Programming Language: Go
- Database: PostgreSQL
- Deployment: Docker, Linux
- Tools: VSCode, Git
- Patterns: DDD, CQRS

## List of APIs

POST /friendship/connect

GET /friendship/friends

GET /friendship/mutuals

POST /subscription/subscribe

POST /subscription/block

GET /subscription/updates_user

## Deployment
This project can be deployed by Docker to Linux server at: http://localhost:3000/
```
make dev
```

Testing
```
make test
```

or run to develop local
```
make setup_db // for setup db
make run_es // run microservice
```

## Layout

```tree
├── ...
├── common/
├── module/
│   └── friendship/
│       ├── adapter/
│       │   └── postgres/
│       │       ├── repository/
│       │       ├── model/
│       │       └── convert/
│       ├── app/
│       │   ├── command/
│       │   │   └── payload/
│       │   └── query/
│       │       └── payload/
│       ├── domain/
│       ├── port/
│       └── service.go
├── mock/
├── middleware/
├── migration/
├── pkg/
├── postman/
├── docker-compose.yml
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── ...
```

A brief description of the friendship module layout:

* `service.go` is the file to inject the repositories and application into port server and router api.
* `port/` is the place convert and validate request from client
* `domain/` is the place hold the core entities business
* `app/` is the place hold the logic business handling
* `adapter/` is the place holder many external technologies like PostgreSQL, Redis, etc