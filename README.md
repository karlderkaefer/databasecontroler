[![Build Status](https://travis-ci.org/karlderkaefer/databasecontroler.png)](https://travis-ci.org/karlderkaefer/databasecontroler)
[![codecov](https://codecov.io/gh/karlderkaefer/databasecontroler/graph/badge.svg)](https://codecov.io/gh/karlderkaefer/databasecontroler)
## Database Controler 

### Docker images for database
oracle image are not available on docker hub because of licensing.
Instead you can build them by own. 
Clone oracle repo and follow build instructions on https://github.com/oracle/docker-images/blob/master/OracleDatabase/SingleInstance/README.md
Start database with `docker-compose up -d oracle11 oracle12`

### Run
```bash
go get github.com/karlderkaefer/databasecontroler
cd $GOPATH/src/github.com/karlderkaefer/databasecontroler
go run main.go
```
