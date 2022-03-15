# How to build and run todo app
* Paste the github CLIENT ID and CLIENT SECRET in the docker-compose.yml file
```sh
docker-compose up todoapp
```
# How to run unit tests
Unit tests are not embedded in the build process but can be run in the docker container with the following command. The unit tests are run against a real postgres instance.
```sh
go test ./...
```

# How to use todo app

Hit `localhost:4000` and login using Github. Complete the login and make a note of the returned JWT token.

# Testing the app

Export the JWT token as env variable for ease of use
```sh
export MY_JWT=<JWT token>
```
## Create some notes

```sh
curl --location --request POST 'localhost:4000/notes' \
--header 'Authorization: Bearer '"$MY_JWT"'' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Sample note 1",
    "body": "This is a sample note 1"
}'
```
```sh
curl --location --request POST 'localhost:4000/notes' \
--header 'Authorization: Bearer '"$MY_JWT"'' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Sample note 2",
    "body": "This is a sample note 2"
}'
```
```sh
curl --location --request POST 'localhost:4000/notes' \
--header 'Authorization: Bearer '"$MY_JWT"'' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Sample note 3",
    "body": "This is a sample note 3"
}'
```

## Read the created notes

```sh
curl --location --request GET 'localhost:4000/notes' \
--header 'Authorization: Bearer '"$MY_JWT"''
```

## Read a particular note

```sh
curl --location --request GET 'localhost:4000/notes/<note-id>' \
--header 'Authorization: Bearer '"$MY_JWT"''
```

## Update a particular note
```sh
curl --location --request PUT 'localhost:4000/notes/<note-id>' \
--header 'Authorization: Bearer '"$MY_JWT"'' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Sample note 10",
    "body": "This is a sample note 10"
}'
```

## Delete a particular note
```sh
curl --location --request DELETE 'localhost:4000/notes/460b4a7d-4662-4a41-a961-596f0d636699' \
--header 'Authorization: Bearer '"$MY_JWT"''
```

## Errors
Only generalised errors are returned in the response. Please check the console logs for the exact errors if any occurred.
