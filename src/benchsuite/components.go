// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package benchsuite

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/glebarez/go-sqlite"

	"regression-ci/internal/config"
	"regression-ci/internal/regression"
	"regression-ci/pkg/types"
)

func benchmarkRegression() []BenchmarkResult {
	var results []BenchmarkResult

	db := setupTestDB()
	defer db.Close()
	detector := setupDetector(db)

	results = append(results, measureBenchmark("regression", "detection", func(b *testing.B) {
		req := createTestRequest()
		for i := 0; i < b.N; i++ {
			detector.Analyze(req)
		}
	}))

	results = append(results, measureBenchmark("regression", "baseline_calculation", func(b *testing.B) {
		values := []float64{95.2, 98.1, 96.8, 97.5, 99.0}
		for i := 0; i < b.N; i++ {
			calculateBaselineMock(values)
		}
	}))

	return results
}

func benchmarkIntegration() []BenchmarkResult {
	var results []BenchmarkResult

	db := setupTestDB()
	defer db.Close()

	results = append(results, measureBenchmark("integration", "database_write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			insertBenchmarkMock(db, "test/repo", "main", "abc123", "component", 100.0)
		}
	}))

	results = append(results, measureBenchmark("integration", "database_read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			queryBaselineMock(db, "test/repo", "component")
		}
	}))

	return results
}

func benchmarkMemory() []BenchmarkResult {
	var results []BenchmarkResult

	results = append(results, measureBenchmark("memory", "json_processing", func(b *testing.B) {
		data := generateLargeBenchmarkJSON()
		for i := 0; i < b.N; i++ {
			parseBenchmarkJSONMock(data)
		}
	}))

	return results
}

func benchmarkConcurrent() []BenchmarkResult {
	var results []BenchmarkResult

	results = append(results, measureBenchmark("concurrent", "parallel_analysis", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				detectRegressionMock(100.0, 90.0, 10.0)
			}
		})
	}))

	return results
}

func detectRegressionMock(current, baseline, threshold float64) bool {
	change := ((current - baseline) / baseline) * 100
	return change > threshold
}

func calculateBaselineMock(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func setupTestDB() *sqlx.DB {
	db, _ := sqlx.Open("sqlite", ":memory:")
	
	// Create tables for testing
	schema := `
	CREATE TABLE IF NOT EXISTS benchmarks (
		id INTEGER PRIMARY KEY,
		repo TEXT NOT NULL,
		branch TEXT NOT NULL,
		commit_hash TEXT NOT NULL,
		component TEXT NOT NULL,
		value REAL NOT NULL,
		timestamp INTEGER NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS baselines (
		repo TEXT NOT NULL,
		component TEXT NOT NULL,
		baseline_value REAL NOT NULL,
		sample_count INTEGER DEFAULT 5,
		updated_at INTEGER NOT NULL,
		PRIMARY KEY (repo, component)
	);`
	
	db.Exec(schema)
	return db
}

func insertBenchmarkMock(db *sqlx.DB, repo, branch, commit, component string, value float64) {
}

func queryBaselineMock(db *sqlx.DB, repo, component string) float64 {
	return 100.0
}

func generateLargeBenchmarkJSON() []byte {
	return []byte(`{"component": "test", "value": 100.0}`)
}

func parseBenchmarkJSONMock(data []byte) map[string]interface{} {
	return map[string]interface{}{"component": "test", "value": 100.0}
}

func setupDetector(db *sqlx.DB) *regression.Detector {
	cfg := config.DetectionConfig{
		DefaultThreshold: 10.0,
		MinSamples:       5,
		MaxSamples:       50,
	}
	return regression.New(db, cfg)
}

func createTestRequest() types.AnalyzeRequest {
	return types.AnalyzeRequest{
		Repo:   "test/repo",
		Branch: "main", 
		Commit: "abc123",
		Components: map[string]float64{
			"test_component": 100.0,
		},
	}
}