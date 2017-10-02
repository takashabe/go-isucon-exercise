# benchmark

## directories

```
bench/
  benchmark/  # CLI benchmark script
  agent/      # Communication with pubsub server and frontend on the benchmarker
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

## agent

run command `agent`. if your will build `make build`.

options:

* interval - polling queues interval
* pubsub - pubsub server URL
* benchmark - benchmark script file path
* param - parameters file path for the benchmark script
* host - webapp hostname
* port - webapp running port

example run command:

```
agent -pubsub=http://localhost:9000 -benchmark=./bin/app -param=./param.json -host=localhost -port=8080
```

### Queue

Polling for the benchmark request queues, and dispatch a request to the benchmarker. The result benchmark send queue when finished benchmark.

Send response queue message:

```
Data: {
  // Benchmark result
}
Attributes: {
  "source_msg_id": // Request message queue ID,
  "team_id":       // Team ID from the request message queue,
  "created_at":    // Benchmark finished time,
},
```
