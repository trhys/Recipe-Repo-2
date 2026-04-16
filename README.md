# The Recipe Repository
---
## Overview

The Recipe Repo is a recipe sharing and shopping helper application. Create/upload recipes with their ingredients and a description and use the shopping list helper to automatically populate a shopping list for your selected recipes.

## Motivation

This idea came to mind originally in order to avoid the hassle of meal planning and trying to think of everything you might need to buy before going to the store. I want to be in and then immediately out when I shop, and that includes the planning stage. Having a list of meals easily selectable and a shopping list created more or less instantly is the goal of this application.[^1] 

[^1]: does not include fast and decisive spouse

## Usage

I intend for this project to be deployed with finished web and android applications. For the purpose of this demonstration however, you can most easily deploy this through Docker(TODO). Otherwise, refer to Quick Start section to deploy manually.

**Note:** This backend does use AWS for S3 and CDN. It's not strictly necessary so some placeholder values won't break it, but recipe images require these to be set up.

---

### Docker

- For maximum convenience, simply ```docker pull timreese/reciperepo:latest``` to pull the docker image.

- To start the container: ```docker run -p 8080:80 reciperepo```

- If you want to set up the AWS variables to use your AWS for image serving: ```docker run -e S3_BUCKET=<your-bucket> -e S3_REGION=<your-region> -e S3_CDN=<your-cdn> -p 8080:80 reciperepo```

**Security note:** This backend server uses JWT authorization for certain endpoints and this demo uses a generic secret in the SECRET env variable. If you plan to use this publicly you should change that.

### Manual Setup

- Clone this repository: ```git clone https://github.com/trhys/Recipe-Repo-2.git```

- You need to setup a POSTGRESQL database server to run this application: [Download it](https://www.postgresql.org/download/)

- Configure your .env variables. See the bottom of this section for an overview of those.

- If your .env is set up and your database is switched on, you can either build (```go build . -o ./reciperepo```) or ```go run .```

If you build the binary, start the server with ```./reciperepo``` assuming you use the same filename as this example.

#### ENV
- DB= your database url string example:```postgres://postgres:postgres@localhost:5432/recipe_repo?sslmode=disable```
- PLATFORM= if this is set to "dev", this enables the /api/reset endpoint which will clear your database.
- SECRET= this is what the auth package will use to validate tokens. keep it secret. keep it safe.
- JWT_DUR= this sets the expiry on authorization tokens. it is in seconds (```3600``` = 1 hour)
- APP_DIR= root path for the frontend. probably doesn't need changed but I won't tell you no
- S3_BUCKET= your s3 bucket for image file storage. you'll need to get this from AWS:```recipe_repo12345```
- S3_REGION= your s3 region for your bucket:```us-east-2```
- S3_CDN= the cdn url from AWS if it's set up

---

If your server is running and everything works correctly, go to ```http://localhost:8080/``` to view the demo web application.
