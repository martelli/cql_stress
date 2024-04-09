package main

import (
	"flag"
	"log"

	"cql_stress/ralpe"
	"cql_stress/randzylla"
)

func main() {

	var (
		runs        int
		parallelism int
		rate        int
		save        bool
		scylladb    string
	)

	flag.IntVar(&rate, "rate-limit", 1, "Number of requests per second")
	flag.IntVar(&parallelism, "parallelism", 1, "Number of parallel workers")
	flag.IntVar(&runs, "runs", 1, "Number of consecutive runs")
	flag.BoolVar(&save, "save", false, "Preserve test data")
	flag.StringVar(&server, "server", "localhost:9042", "ScyllaDB IP:port")

	flag.Parse()

	rz, err := randzylla.NewRandzylla(server)
	if err != nil {
		log.Fatal(err)
	}

	insertFunction := rz.GetInsertFunction()

	r := ralpe.NewRalpe(insertFunction, rate, parallelism, rate*runs)

	r.Start()

	r.Wait()

	if !save {
		err = rz.TearDown()
		if err != nil {
			log.Fatal(err)
		}
	}
}
