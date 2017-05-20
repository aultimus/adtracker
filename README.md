# adtracker

adtracker is a persistent RESTful ad tracking service

## Installation ##
```bash
go get -u github.com/aultimus/adtracker/adtracker
```

## Usage ##

start redis
```bash
$ redis-server
```

start adtracker server
```bash
$ adtracker
```

ad_count and track endpoints serve as get and increment operations for the given ID
example client usage:

```bash
$ curl localhost:5000/ad_count -X GET --data '{"ID":"foo"}'
ID foo does not exist
$ curl localhost:5000/track -X PUT --data '{"ID":"foo"}'
$ curl localhost:5000/ad_count -X GET --data '{"ID":"foo"}'
{"Value":1}
$ curl locatrack -X PUT --data '{"ID":"foo"}'
$ curl localhost:5000/ad_count -X GET --data '{"ID":"foo"}'
{"Value":2}
```
