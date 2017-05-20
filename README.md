# adtracker

adtracker is a RESTful persistent key-value store

## Usage ##

start redis
```bash
$ redis-server
```

ad_count and track endpoints serve as get and increment operations

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
