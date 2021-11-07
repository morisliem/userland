# Userland

This repository contains impelementation of "Userland" on boarding project

Userland is an imaginary authentication and session tracking service that is defined in [this Apiary ](https://userland.docs.apiary.io)

This implementation is going to have extra requirements:

- Password minimum 8 characters, has lowercase, uppercase, number

- Forgot password: must be different from last 3 passwords

- Use at some common JWT tokens payload

- Use JTI to revoke session https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7 

- Use prepared statement

- Change email, needs to send confirmation email to the new email

- Delete account = soft delete

- Client list can be added on demand with  “upsert”

- Password is hashed on rest, hash using bcrypt

- OTP has timeout and immediately revoke OTP when it is used

  

## API Contract

https://userland.docs.apiary.io/#introduction/http-status-codes

## How to run the code

### Step 1: Clean download libraries

from the root folder, run `go mod tidy`



### Step 2: Environment setup

run `docker-compose up` from the root repository

Postgres image is from https://hub.docker.com/_/postgres

By default, the internal database is served in `db_userland:5432` , with connection string: `postgres://admin:admin@db_userland:5432/userland` Postgres data is stored in `{root_repo}/data`. Those defaults can be modified in:

```
postgres:
    image: postgres:14-alpine
    restart: on-failure
    volumes:
        - ./data:/var/lib/postgresql/data # default postgres data location
    environment:
        - POSTGRES_USER=admin # default username
        - POSTGRES_PASSWORD=admin # default password
        - POSTGRES_DB=userland # default db name
        - PGPORT=5432 # default port
    networks: 
      default:
        aliases: 
          - db_userland # default host
    expose:
      - 5432 # default port exposed in docker default network
```

To inspect the database with GUI, use the `adminer` (from https://hub.docker.com/_/adminer) via `0.0.0.0:8081`. The default port can be configured in

```
adminer:
    ...
    ports:
      - "8081:8080/tcp" # default to 8081
		...
```



### Step 3: Using the app

By default, the service is running on `:8080` , that can be modified by changing the docker compose file

```
userland:
	...
	ports:
  	- "8080:80/tcp" # set 8080 to other port
  ...
```



