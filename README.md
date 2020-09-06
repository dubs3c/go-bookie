<h1 align="center">Go-Bookie</h1>
<p align="center">
<a href="https://goreportcard.com/badge/github.com/dubs3c/go-bookie"><img src="https://goreportcard.com/badge/github.com/dubs3c/go-bookie"></a>
<a href="https://img.shields.io/github/workflow/status/dubs3c/go-bookie/Go"><img src="https://img.shields.io/github/workflow/status/dubs3c/go-bookie/Go"></a>
</p>

<p align="center">
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