# social-hub

## Description

A social network with the following features implemented:

- followers/following
- profiles
- posts
- groups
- notifications
- private and group chats

## Running the project locally

## Running the backend server

```console
go run ./api/.
```

If you wish to seed the database, run the command:

```console
go run ./api/. seed
```

## Running the frontend server

```console
$ cd frontend
$ npm install
$ npm start
```

## Running the project with Docker

### Run Docker scripts

```console
docker compose up -d
```

### Clean up Docker files

```console
bash docker-remove.sh
```

## Database

- Creating new database migration:

```console
migrate create -ext sql -dir api/pkg/db/migrations/sqlite -seq schema_name
```

"schema_name" is the name of the migration

- Updating the migration.go file code:

```console
go-bindata -o api/pkg/db/sqlite/migration.go -prefix api/pkg/db/migrations/sqlite/ -pkg database api/pkg/db/migrations/sqlite
```

## Stack

Frontend

- React
- HTML
- CSS

Backend

- Go
- SQLite3
