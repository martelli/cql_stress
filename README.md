## Design

`ralpe` (rate limited parallel executor) allows you to run any
function of type `func() error` in parallel batches.

`randzylla` will create a keyspace in the ScyllaDB and provide
a function to insert a random into the database. 

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

## Notes

When running many requests in random mode, the cql driver will overwrite entries with
the same value. This can cause the table to display a smaller number of entries.
To assess the exact amount of entries is being written, use the `-serial` flag.
This assumes also the use of the `-save` flag to preverse database content.

## Resource limiting and result sample

To limit the amount of CPU used on both server and client, we can do:

- start scylladb instance we can add `--smp 1`.
- do `export GOMAXPROCS=1` before calling the cli tool.

By using these on AMD Ryzen 7 3700X, we verified average throughput of 100k writes/s
with an average latency around 1ms, while CPU stays around 80% utilization for each
process.

## Help

```
Usage of ./cql_stress:
  -parallelism int
    	Number of parallel workers (default 1)
  -rate-limit int
    	Number of requests per second (default 1)
  -runs int
    	Number of consecutive runs (default 1)
  -save
    	Preserve test data
  -serial
    	Use serial values instead of random generated
  -server string
    	ScyllaDB IP:port (default "localhost:9042")
```
