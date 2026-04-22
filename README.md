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
- PLATFORM= if this is set to "dev", this enables the /api/reset endpoint which will clear your database.
- SECRET= this is what the auth package will use to validate tokens. keep it secret. keep it safe.
- JWT_DUR= this sets the expiry on authorization tokens. it is in seconds (```3600``` = 1 hour)
- APP_DIR= root path for the frontend. probably doesn't need changed but I won't tell you no
- ADMIN_DIR= path to the admin folder in /app. Separate from the base fileserver
- S3_BUCKET= your s3 bucket for image file storage. you'll need to get this from AWS
- S3_REGION= your s3 region for your bucket
- S3_CDN= the cdn url from AWS if it's set up
- IMAGE_PLACEHOLDER= a key that pulls an image from S3 - use for recipes that have no image uploaded

**Security note:** This backend server uses JWT authorization for certain endpoints and this demo uses a generic secret in the SECRET env variable. If you plan to use this publicly you should change that.

#### Docker

- If your env is set up correctly, all you need to do now is ```docker compose up``` and the containers will spin up and bind to localhost:8080

- If the db url is not set correctly, the container will start with no database connection.

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

---

If your server is running and everything works correctly, go to ```http://localhost:8080/``` to view the demo web application.

## API DOCUMENTATION

**A note on response shapes:** because most of the functionality uses the same basic shape, there are only a couple structs that shape the JSON data and some is omitted case to case based on the endpoint.

#### User endpoints:

##### Basic response shape for endpoints that return user data:

```
  {
     "id" : a uuid,
     "created_at" : a timestamp when the user was created.
     "updated_at" : last updated timestamp of user,
     "email" : user email address,
     "name" : user display name,
  }
  ```
---

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
---

#### Recipe Endpoints

##### Basic recipe response shape: 
```
{
  "id": uuid,
  "title": title of recipe,
  "created_at": timestamp of creation,
  "updated_at": timestamp of last update,
  "user_id": uuid of creator,
  "author": display name of creator,
  "description": the description,
  "image_key": key is used to search s3 for image asset,
  "ingredients": array of ingredient objects,
  "image_url": constructed url for image,
}
```

Some responses give ingredients in this shape:
```
{
  "id": ingredient uuid,
  "name": display name,
  "quantity": float32 quantity value,
  "unit": string unit measurment,
}
```
---

- "GET /users/{user_id}" : takes header Accept - if set to application/json, responds with json in this shape:
  ```
  {
    "recipes": array of recipes,
    "name": display name for {user_id},
  }
  ```
  If Accept != application/json, this endpoint renders an HTML template.
  
- "GET /recipes/{recipe_id}" : takes header Accept - if set to application/json, responds with json in the basic recipe shape. If Accept != application/json, this endpoint renders an HTML template.
  
- "GET /api/recipes" : responds with ten most recently created recipes from the database. The json shape is:
  ```
  {
    "recipes": array of recipes in the basic shape,
  }
  ```

- "POST /api/new_recipe" : requires Authorization header (Authorization: Bearer $jwt_token). Takes a request body and creates the recipe in the database. Returns the recipe in the base json shape above.
  REQUEST:
  ```
  {
    "title": recipe title,
    "user_id": creator's id,
    "description": description provided by user,
    "ingredients": array of ingredients,
  }
  ```
---

#### Ingredient Endpoints

##### Basic ingredient response shape: 
```
{
   "id': uuid,
   "name": display name,
   "image_key": key to aws asset,
   "created_at": timestamp of creation,
   "updated_at": timestamp of last update,
}
```
---

- "POST /api/admin/new_ingredient" - create new ingredient. currently an admin only feature but will change. 'admin' is a column in the users table that is checked here
  REQUEST:
  ```
  {
     "name": display name for the ingredient,
  }
  ```
  This ep will take more in the future as other API integrations occur.

- "GET /api/get_ingredients" - responds with all the ingredients available in the database. Will probably make some changes as this is likely not scalable.
  RESPONSE:
  ```
  {
     "ingredients": array of basic ingredient responses,
  }
  ```
---

#### Shopping list Endpoints

##### Basic shopping list response shape: 
```
{
   "id": uuid,
   "name": display name,
   "created_at": timestamp of creation,
   "updated_at": timestamp of last update,
}
```
---

- "GET /shoppinglists/{shopping_list_id} - takes "Accept" header: if application/json:
  RESPONSE:
  ```
  {
     "id": uuid,
	  "name": display name,
     "created_at": timestamp of creation,
     "recipes": array of recipe responses,
     "quantity": map of recipes to quantity selected,
  }
  ```
  Otherwise renders HTML template.
  
- "GET /users/{user_id}/shoppinglists" - gets shopping lists for user. takes header: Authorization Bearer $token, && "Accept"
  if Accept: application/json:
  ```
  {
     "name": user's display name,
     "shopping_lists": array of shopping list responses,
  }
  ```
  Otherwise render HTML
  
- "POST /api/new_shopping_list" - create new shopping list. takes Authorization Bearer token header
  REQUEST:
  ```
  {
     "name": shopping list name,
  }
  ```
  Responds with created list
  
- "POST /api/add_to_list" - add recipe to list
  REQUEST:
  ```
  {
     "shopping_list_id": uuid of shopping list,
     "recipe_id": uuid of recipe,
     "quantity": number to add,
  }
  ```
  No response body given
---

#### Token Endpoints

- "POST /api/tokens/refresh" - takes Authorization Bearer token header, where token is a refresh token. Checks token in db and returns new JWT if valid
  RESPONSE:
  ```
  {
     "token": JWT,
  }
  ```
  
- "POST /api/tokens/revoke" - takes Authorization Bearer token, where token is a refresh token. Returns no body, just revokes existing token and respond with status 204

## Contributing

todo
