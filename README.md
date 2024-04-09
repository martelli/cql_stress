## Execution

1. Run ScyllaDB under docker:

```
docker run -it -p 9042:9042 scylladb/scylla
```

2. Checkout and build the `cql_stress` tool:

```
mkdir proj && cd proj
git clone https://github.com/martelli/cql_stress
cd cql_stress
go build
```

3. Run it against the ScyllaDB server:

```
./cql_stress -parallelism 100 -rate-limit 100000 -runs 5 -server 192.168.1.10:9042
2024/04/09 18:16:35 INFO Stats: throughput=100000 avg_latency:=696.249µs
2024/04/09 18:16:36 INFO Stats: throughput=100000 avg_latency:=728.238µs
2024/04/09 18:16:37 INFO Stats: throughput=100000 avg_latency:=706.838µs
2024/04/09 18:16:38 INFO Stats: throughput=100000 avg_latency:=716.805µs
2024/04/09 18:16:39 INFO Stats: throughput=100000 avg_latency:=727.427µs
```

## Help

Usage of ./cql_stress:
  -parallelism int
    	Number of parallel workers (default 1)
  -rate-limit int
    	Number of requests per second (default 1)
  -runs int
    	Number of consecutive runs (default 1)
  -save
    	Preserve test data
  -server string
    	ScyllaDB IP:port (default "localhost:9042")
