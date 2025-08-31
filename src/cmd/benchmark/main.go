// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package main

import (
	"flag"
	"fmt"
	"os"

	"regression-ci/benchsuite"
)

func main() {
	outputFile := flag.String("o", "benchmark_results.json", "Output file for benchmark results")
	flag.Parse()

	fmt.Println("Running benchmark suite...")
	
	suite, err := benchsuite.RunSuite()
	if err != nil {
		fmt.Printf("Error running benchmark suite: %v\n", err)
		os.Exit(1)
	}

	if err := benchsuite.WriteResults(suite, *outputFile); err != nil {
		fmt.Printf("Error writing benchmark results: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Benchmark suite completed. Results written to %s\n", *outputFile)
}