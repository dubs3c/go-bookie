# go-bookie

A new version of bookie, written in Go :)

**Run docker container**
```
docker run -e POSTGRES_PASSWORD=bookie -e POSTGRES_USER=bookie -e POSTGRES_DB=bookie -p 5432:5432 -d postgres
```

**Run Migrations**
```
migrate -database postgres://bookie:bookie@localhost:5432/bookie?sslmode=disable -path migrations up
```