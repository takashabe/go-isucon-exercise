# benchmark

## directories

```
bench/
  benchmark/  # CLI benchmark script
  server/     # call benchmark HTTP server
```

#### benchmark

run command `bin/app`. if your will build `make build`.

options:

* port - choose webapp server port
* host - choose webapp server host
* file - path of webapp application account list file
* agent - used benchmark request user agent

example:

```
bin/app -host=localhost -port=8080 -file=data/param.json -agent=isucon_go
```

#### server

listen and serve benchmark queue request.

run benchmark request URL:

```
/api/benchmark/:entry_id
```

client need enqueue at `:entry_id` before request. server receive request, dequeue entry_id resource and update queue.
