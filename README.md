# ConnectTeam server app 

## Run 
1. First, install go and golang-migrate. For macOS:
``` bash
brew install golang-migrate
```
2. Add .env file in the project.
3. Run docker container. Examle:
``` bash
docker run --name=ct-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
```
4. Run migrations:
``` bash
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
```
5. Run app:
``` bash
go run cmd/main.go
```
