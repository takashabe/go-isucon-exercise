# benchmark

## directories

```
bench/
  benchmark/  # CLI benchmark script
  server/     # call benchmark HTTP server
```

## benchmark

run command `bin/app`. if your will build `make build`.

options:

* port - choose webapp server port
* host - choose webapp server host
* file - path of webapp application account list file
* agent - used benchmark request user agent

example run command:

```
bin/app -host=localhost -port=8080 -file=data/param.json -agent=isucon_go
```

#### result

Benchmark returns JSON result that include fields `valid`, `request_count`, `elapsed_time` and so.

example benchmark result when succeed:

```
{
        "valid": true,
        "request_count": 3651,
        "elapsed_time": 0,
        "response": {
                "success": 1452,
                "redirect": 2199,
                "client_error": 0,
                "server_error": 0,
                "exception": 0
        },
        "violations": []
}
```

when failed benchmark:

```
{
        "valid": false,
        "request_count": 1,
        "elapsed_time": 0,
        "response": {
                "success": 0,
                "redirect": 0,
                "client_error": 0,
                "server_error": 1,
                "exception": 0
        },
        "violations": [
                {
                        "request_type": "INITIALIZE",
                        "description": "パス '/initialize' からレスポンスが返ってきませんでした",
                        "num": 1
                }
        ]
}
```

## server

listen and serve benchmark queue request.

run benchmark request URL:

```
/api/benchmark/:entry_id
```

client need enqueue at `:entry_id` before request. server receive request, dequeue entry_id resource and update queue.
