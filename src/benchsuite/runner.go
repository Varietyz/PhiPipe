package benchsuite

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

type BenchmarkResult struct {
	Component   string        `json:"component"`
	Operation   string        `json:"operation"`
	Duration    time.Duration `json:"duration_ns"`
	Iterations  int           `json:"iterations"`
	Memory      int64         `json:"memory_bytes"`
	Allocations int64         `json:"allocations"`
	Timestamp   int64         `json:"timestamp"`
}

type SuiteResult struct {
	Results   []BenchmarkResult `json:"results"`
	GoVersion string            `json:"go_version"`
	GOOS      string            `json:"goos"`
	GOARCH    string            `json:"goarch"`
	Timestamp int64             `json:"timestamp"`
}

func RunSuite() (*SuiteResult, error) {
	var results []BenchmarkResult

	regression := benchmarkRegression()
	integration := benchmarkIntegration()
	memory := benchmarkMemory()
	concurrent := benchmarkConcurrent()

	results = append(results, regression...)
	results = append(results, integration...)
	results = append(results, memory...)
	results = append(results, concurrent...)

	suite := &SuiteResult{
		Results:   results,
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		Timestamp: time.Now().Unix(),
	}

	return suite, nil
}

func WriteResults(suite *SuiteResult, filename string) error {
	data, err := json.MarshalIndent(suite, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write results: %w", err)
	}

	return nil
}

func measureBenchmark(name, operation string, fn func(*testing.B)) BenchmarkResult {
	result := testing.Benchmark(fn)

	return BenchmarkResult{
		Component:   name,
		Operation:   operation,
		Duration:    time.Duration(result.NsPerOp()),
		Iterations:  result.N,
		Memory:      int64(result.MemBytes),
		Allocations: int64(result.MemAllocs),
		Timestamp:   time.Now().Unix(),
	}
}