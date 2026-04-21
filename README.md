# The Recipe Repository

## Overview

A Golang backend with a Postgresql database, containerized with Docker Compose: *The Recipe Repo* is a recipe sharing and shopping helper application. Create/upload recipes with their ingredients and a description and use the shopping list helper to automatically populate a shopping list for your selected recipes.

## Motivation

This idea came to mind originally in order to avoid the hassle of meal planning and trying to think of everything you might need to buy before going to the store. I want to be in and then immediately out when I shop, and that includes the planning stage. Having a list of meals easily selectable and a shopping list created more or less instantly is the goal of this application.[^1] 

[^1]: does not include fast and decisive spouse

## Usage

I intend for this project to be deployed with finished web and android applications. For the purpose of this demonstration however, you can most easily deploy this with Docker Compose.

**Note:** This backend does use AWS for S3 and CDN. It's not strictly necessary so some placeholder values won't break it, but recipe images require these to be set up.

---

- Clone this repository: ```git clone https://github.com/trhys/Recipe-Repo-2.git```

#### Environment

Certain environment variables are required to configure the server. You can edit the .env-example and rename the file: ```mv .env-example .env```

Note: the database url string will differ for running the server in a docker container. 

For localhost: ```postgres://postgres:postgres@localhost:5432/recipe_repo?sslmode=disable```

For docker: ```postgres://postgres:postgres@db:5432/recipe_repo?sslmode=disable```

- DB= your database url string example:```postgres://postgres:postgres@localhost:5432/recipe_repo?sslmode=disable```
- PLATFORM= if this is set to "dev", this enables the /api/reset endpoint which will clear your database.
- SECRET= this is what the auth package will use to validate tokens. keep it secret. keep it safe.
- JWT_DUR= this sets the expiry on authorization tokens. it is in seconds (```3600``` = 1 hour)
- APP_DIR= root path for the frontend. probably doesn't need changed but I won't tell you no
- S3_BUCKET= your s3 bucket for image file storage. you'll need to get this from AWS:```recipe_repo12345```
- S3_REGION= your s3 region for your bucket:```us-east-2```
- S3_CDN= the cdn url from AWS if it's set up
- IMAGE_PLACEHOLDER= a key that pulls an image from S3 - use for recipes that have no image uploaded

**Security note:** This backend server uses JWT authorization for certain endpoints and this demo uses a generic secret in the SECRET env variable. If you plan to use this publicly you should change that.

#### Docker

- If your env is set up correctly, all you need to do now is ```docker compose up``` and the containers will spin up and bind to localhost:8080

- If the db url is not set correctly, the container will start with no database connection.

#### Manual Setup

- You need to setup a POSTGRESQL database server to run this application: [Download it](https://www.postgresql.org/download/)

- Configure your .env variables. See the bottom of this section for an overview of those.

- If your .env is set up and your database is switched on, you can either build (```go build . -o ./reciperepo```) or ```go run .```

If you build the binary, start the server with ```./reciperepo``` assuming you use the same filename as this example.

---

If your server is running and everything works correctly, go to ```http://localhost:8080/``` to view the demo web application.

## API

#### User endpoints:

Basic response shape for endpoints that return user data:

```
  {
     "id" : a uuid,
     "created_at" : a timestamp when the user was created.
     "updated_at" : last updated timestamp of user,
     "email" : user email address,
     "name" : user display name,
  }
  ```

- "GET /api/users/{user_id}" : responds with JSON user information

- "POST /api/new_user" : creates a new user in the database. takes an email, username, and password in this shape:

  REQUEST:
  ```
  {
    "email" : an email address,
    "name" : user display name,
    "password": user password,
  }
  ```

  The backend will hash the password and store the hash in the users table.
  
- "POST /api/login" : send request with email and password and the backend will check against the hash stored in the users table and respond with a JWT and refresh token if validated.

  REQUEST:
  ```
  {
    "email": email address,
    "password": user's password,
  }
  ```

  RESPONSE:
  ```
  {
    "id": user's uuid,
    "email": user's email address",
    "username": display name,
    "token": JWT for authorization(1 hr expiry by default),
    "refresh_token": token stored in database with 30 day expiry,
  }
  ```
  
- "POST /api/admin/reset" : takes no request body, just requires the PLATFORM env variable to be set to "dev". Empties the users table which will cascade onto mostly everything else. Use this to reset the database if desired.
