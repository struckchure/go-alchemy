# Alchemy (Go)

Read to use components for backend developers.

# Philosophy (is that spelling correct? 🤔)

If frontend developers can do `npx shadcn add button`, backend developers should also be able to do `alchemy add authentication`.

# Installation

Clone the repo and build with this command;

```sh
$ go build -o alchemy ./cmd
```

# Usage

```sh
alchemy-sandbox % alchemy init
? Provide alchemy component root:  .
? Choose ORM:  Prisma
? Choose Database Provider:  PostgreSQL
? Provision Database with Docker Compose:  Yes
Creating PostgreSQL with Docker Compose
Docker Compose file successfully generated
Using Prisma ORM
Downloading go prisma client
Initializing new prisma project
✨ Alchemy config has been generated!
🛠️  Updating Go dependencies ...
🥂 You're all set!

Start Database Service
$ docker compose up -d


Interactively add component
$ alchemy add Authentication // this will add all components from the authentication module

Or add a specific component
$ alchemy add Authentication.Login

alchemy-sandbox % alchemy add
? Component Category Authentication
? Select Components Login
Creating Authentication.Login component
  + prisma/schema.prisma
  + dao/user_dao.go
  + services/authentication_service.go
+ Authentication.Login
```
