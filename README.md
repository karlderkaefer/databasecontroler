[![Build Status](https://travis-ci.org/karlderkaefer/databasemanager.png)](https://travis-ci.org/karlderkaefer/databasemanager)
[![codecov](https://codecov.io/gh/karlderkaefer/databasemanager/branch/master/graph/badge.svg)](https://codecov.io/gh/karlderkaefer/databasemanager)
[![Go Report Card](https://goreportcard.com/badge/karlderkaefer/databasemanager)](https://goreportcard.com/report/github.com/karlderkaefer/databasemanager)
## Database Manager

### Oracle
oracle image are not available on docker hub because of licensing.
Instead you can build them by own. 
Clone oracle repo and follow build instructions on https://github.com/oracle/docker-images/blob/master/OracleDatabase/SingleInstance/README.md
Start database with `docker-compose up -d oracle11 oracle12`

### Run
```bash
go get github.com/karlderkaefer/databasemanager
cd $GOPATH/src/github.com/karlderkaefer/databasemanager
go run main.go
# in developing mode you can use fresh
go get github.com/pilu/fresh
fresh
```

### DB2
db2 is a pain here. db2 does not allow you to create databases with jdbc.
Those actions have to be done with db2 command line. 
Installing db2cli is a painful as well.
To workaround this we are using existing docker image and create same aliases.
Linux is straightforward but for Windows a little bit tricky.
Add `db2.cmd` to `%USERPROFILE%\bin` and add this directory to your `PATH`.
Content of `db2.cmd`
```bash
@echo off
docker exec databasemanager_db2_1 su - db2inst1 -c "/home/db2inst1/sqllib/bin/db2 %*"
```
Now you are able to run db2 command directly from host in the container.
```bash
db2 create database test
```
I dont want to run the database manager on same host as db2 is running.
Therefor we simulate a db2 remote database with the container `db2remote`.
`db2` container will connect to remote server and execute commands remotely.
Note that database test has to exist on `db2remote`. Prepare the `db2` container:
```bash
db2 catalog tcpip node remoteDb2 remote db2remote server 50000
db2 terminate
db2 list node directory 
db2 attach to test user db2inst1 using HelloApes66
db2 create database test 
```
