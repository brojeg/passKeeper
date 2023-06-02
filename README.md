# passKeeper



# Sever Overview
The passKeeper Server is a robust tool designed for the secure handling and management of secrets. It interacts with a database to store and retrieve secrets, uses JWT for authorization, and has functionalities to create tables and define JWT configuration.

### Setup
Before running the server, make sure you have the necessary environment variables set. If these variables are not set, you can alternatively use flags or  .env file while running the application. Here are the environment variables required:

```
RUN_ADDRESS : The address at which the server will run (default is 127.0.0.1:8080).
DATABASE_URI : The connection string for your PostgreSQL database.
JWT_PASSWORD : The password used for JWT.
EXPIRATION_TIME : The TTL (Time To Live) for the JWT token in minutes (default is 15).
```

You can use the following flags in place of environment variables:

```
-a to set server address
-d to set database connection string
-p to set JWT password
-t to set JWT token TTL
```
### Features
HTTP Server: The main server that handles all incoming requests.
External Dependency: Currently, it's the database connection string required to connect to a PostgreSQL database.
Server Auth: Handles server authentication using JWT.



# Client Overview
passKeeper is a robust tool that allows for secure handling and management of secrets. This includes generating new secrets, editing existing ones, listing all stored secrets, and even deleting them when no longer needed.

## Client Commands

### Setup
Sets up initial configurations for passKeeper. This includes setting up the username and password.
`passKeeper setup`


### Login
Initiates login process for a user with an existing passKeeper account.
`passKeeper login`


### Logout
Clears all locally stored passKeeper configuration.
`passKeeper logout`


### New
Generate a new secret of a specific type. Options include key-value pair (kv), credit card details (cc), text (txt), or file.
`passKeeper new [txt|file|kv|cc]`


### List
Displays a list of all secrets currently stored in passKeeper.
`passKeeper list`


### Delete
Removes a secret stored in passKeeper by its unique identifier.
`passKeeper delete [secret_id]`


### Edit
Edits the contents of a secret stored in passKeeper by its unique identifier.
`passKeeper edit [secret_id]`


### Describe
Provides comprehensive details of a secret stored in passKeeper by its unique identifier.
`passKeeper describe [secret_id]`


### Dump
Extracts and exports the binary data of a secret by its unique identifier on the disk.
`passKeeper dump [secret_id]`