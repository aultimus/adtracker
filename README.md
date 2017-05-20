# adtracker

adtracker is a RESTful persistent key-value store

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

ad_count and track endpoints serve as get and increment operations.
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
