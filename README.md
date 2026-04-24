# Overview

A Golang backend with a Postgresql database, containerized with Docker Compose: *The Recipe Repo* is a recipe sharing and shopping helper application. Create/upload recipes with their ingredients and a description and use the shopping list helper to automatically populate a shopping list for your selected recipes.

<img width="720" height="500" alt="Screenshot of a recipe on The Recipe Repository" src="https://github.com/user-attachments/assets/ae5cdade-9757-4931-b917-65996a8ef9d5" />

## Motivation

This idea came to mind originally in order to avoid the hassle of meal planning and trying to think of everything you might need to buy before going to the store. I want to be in and then immediately out when I shop, and that includes the planning stage. Having a list of meals easily selectable and a shopping list created more or less instantly is the goal of this application.[^1] 

[^1]: does not include fast and decisive spouse

## Usage

I intend for this project to be deployed with finished web and android applications. For the purpose of this demonstration however, you can most easily deploy this with Docker Compose. The web demo is viewable in browser at ```http://localhost:8080/```

This readme assumes for the most part you are running on Unix or WSL.

**Note:** This backend does use AWS for S3 and CDN. It's not strictly necessary so some placeholder values won't break it, but recipe images require these to be set up.

#### Prerequisites

- Git
- Docker (recommended)
- Go (if not using docker)

## Quick Start

1. ```git clone https://github.com/trhys/Recipe-Repo-2.git```
2. ```docker compose up```
3. ???
4. profit

In reality, run docker compose from the cloned root. Docker will automatically start the container and bind port 8080. See the demo web app @ http://localhost:8080/
   
---

## Longer Start

- Clone this repository: ```git clone https://github.com/trhys/Recipe-Repo-2.git```

#### Environment

<img width="521" height="186" alt="Screenshot of .env-example" src="https://github.com/user-attachments/assets/a645b258-18f4-47c4-886b-c5faca06ed22" />

---

Certain environment variables are required to configure the server. You can edit the .env-example and rename the file: ```mv .env-example .env``` but you should be able to just build this without changing anything if you don't want to set up the AWS functionality.

Note: the database url string will differ for running the server in a docker container. 

For localhost: ```postgres://postgres:postgres@localhost:5432/recipe_repo?sslmode=disable```

For docker: ```postgres://postgres:postgres@db:5432/recipe_repo?sslmode=disable```

- DB= your database url string
- PLATFORM= no longer used but may return
- SECRET= this is what the auth package will use to validate tokens. keep it secret. keep it safe.
- JWT_DUR= this sets the expiry on authorization tokens. it is in seconds (```3600``` = 1 hour)
- APP_DIR= root path for the frontend. probably doesn't need changed but I won't tell you no
- ADMIN_DIR= path to the admin folder in /app. Separate from the base fileserver
- S3_BUCKET= your s3 bucket for image file storage. you'll need to get this from AWS
- S3_REGION= your s3 region for your bucket
- S3_CDN= the cdn url from AWS if it's set up
- IMAGE_PLACEHOLDER= a key that pulls an image from S3 - use for recipes that have no image uploaded

**Security note:** This backend server uses JWT authorization for certain endpoints and this demo uses a generic secret in the SECRET env variable. If you plan to use this publicly you should change that.

<br>

#### Docker

- If your env is set up correctly, all you need to do now is ```docker compose up``` and the containers will spin up and bind to localhost:8080

- If the db url is not set correctly, the container will start with no database connection.

<br>

#### Manual Setup

- You need to setup a POSTGRESQL database server to run this application: [Download it](https://www.postgresql.org/download/)

- You might need to start the postgresql service ```sudo service postgresql start```

- You'll have to manually set up the database. ```psql -U postgres -h localhost -p 5432 -c "CREATE DATABASE recipe_repo"```
  It will prompt for a password, you can use "postgres" if you don't want to change the db url in the .env file. Otherwise,
  you can use whatever user, password, and database name so long as you set the DB variable correctly.

- [Configure your .env variables](#environment).

- Build the binary: ```go build . -o ./reciperepo```

- Migrations: this application is set up to use [Goose](https://github.com/pressly/goose) for database migrations. Assuming your env is set up and the binary is built, the easiest is to run the entrypoint script ```./entrypoint.sh``` You may need to modify permissions to run it ```chmod +x entrypoint.sh```

- If the script fails or you prefer to run it manually: ```goose -dir sql/schema postgres $your_db_url_string up```

- Start the server ```./reciperepo```. It will seed the database with ingredients from the setup.json file.

<br>

If your server is running and everything works correctly, go to ```http://localhost:8080/``` to view the demo web application.

<br><br>


## API DOCUMENTATION

See full API docs on the [wiki](https://github.com/trhys/Recipe-Repo-2/wiki)

Response bodies take the form of either a JSON body or rendered HTML, depending on the Accept header on certain endpoints.

The [/internal/viewmodel/](https://github.com/trhys/Recipe-Repo-2/tree/rest-and-dry-refactor/internal/viewmodel) package defines the shape of each response.

## Contributing

todo
