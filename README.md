<h1 align="center">Go-Bookie</h1>
<p align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/dubs3c/go-bookie)](https://goreportcard.com/report/github.com/dubs3c/go-bookie) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/dubs3c/go-bookie/Go)

A new version of bookie, written in Go :)

</p>


**Run docker container**
```
docker run -e POSTGRES_PASSWORD=bookie -e POSTGRES_USER=bookie -e POSTGRES_DB=bookie -p 5432:5432 -d postgres
```

**Run Migrations**
```
migrate -database postgres://bookie:bookie@localhost:5432/bookie?sslmode=disable -path migrations up
```