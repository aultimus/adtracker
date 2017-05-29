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
$ curl localhost:5000/ad_count/foo
ID foo does not exist
$ curl localhost:5000/track/foo -X PUT
$ curl localhost:5000/ad_count/foo
{"Value":1}
$ curl localhost:5000/track/foo -X PUT
$ curl localhost:5000/ad_count/foo
{"Value":2}
```

## Design Discussion ##

This service is implemented in golang as it is well suited to the webservice problem domain. Everything needed to run the http server is available in the standard library. It is also very easy to distribute, as a static binary easily obtainable. It is also more performant for a real time system under load than other choices such as python.

Redis was used as a datastore as it is in memory and therefore performant but also satisfies the persistence criteria. The simplicity of the data and queries meant that a relational database would have been uneccesary and also harder to scale.

It would be possible to scale this service horizontally onto multiple instances using a load balancer to balance traffic between adtracker instances. The multiple instances could use a single redis instance, if redis becomes a bottleneck then redis partitioning could be investigated in order to increase throughput, however this would only be necessary with extermely heavy usage.
