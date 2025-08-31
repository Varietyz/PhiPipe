// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/glebarez/go-sqlite"

	"regression-ci/internal/config"
	"regression-ci/internal/server"
	"regression-ci/pkg/types"
)

type TestEnvironment struct {
	db     *sqlx.DB
	server *httptest.Server
	client *http.Client
}

func SetupTestEnv(t *testing.T) *TestEnvironment {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Address:      ":0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: config.DatabaseConfig{
			Path: ":memory:",
		},
		Detection: config.DetectionConfig{
			DefaultThreshold: 10.0,
			MinSamples:       5,
			MaxSamples:       50,
		},
	}

	srv := server.New(db, cfg)
	testServer := httptest.NewServer(srv.Router())

	return &TestEnvironment{
		db:     db,
		server: testServer,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (env *TestEnvironment) Cleanup() {
	env.server.Close()
	env.db.Close()
}

func TestFullPipeline(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("health_check", func(t *testing.T) {
		resp, err := env.client.Get(env.server.URL + "/health")
		if err != nil {
			t.Fatalf("health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("analyze_endpoint", func(t *testing.T) {
		benchmarkData := loadTestBenchmark(t, "golang-project/benchmarks.json")
		
		payload, err := json.Marshal(benchmarkData)
		if err != nil {
			t.Fatalf("failed to marshal benchmark data: %v", err)
		}

		resp, err := env.client.Post(
			env.server.URL+"/analyze",
			"application/json",
			bytes.NewBuffer(payload),
		)
		if err != nil {
			t.Fatalf("analyze request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotImplemented {
			t.Skip("analyze endpoint not yet implemented")
		}
	})

	t.Run("webhook_simulation", func(t *testing.T) {
		webhookPayload := loadWebhookPayload(t, "pr_opened.json")
		
		resp, err := env.client.Post(
			env.server.URL+"/webhook",
			"application/json",
			bytes.NewBuffer(webhookPayload),
		)
		if err != nil {
			t.Fatalf("webhook request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotImplemented {
			t.Skip("webhook endpoint not yet implemented")
		}
	})
}

func TestRegressionScenarios(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	scenarios := []string{"regression.json", "baseline.json", "improvements.json"}
	
	for _, scenario := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			scenarioData := loadBenchmarkScenario(t, scenario)
			validateScenario(t, env, scenarioData)
		})
	}
}

func TestConcurrentLoad(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	concurrency := 10
	requests := 50
	
	results := make(chan error, concurrency*requests)
	
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requests; j++ {
				resp, err := env.client.Get(env.server.URL + "/health")
				if err != nil {
					results <- fmt.Errorf("request failed: %w", err)
					continue
				}
				resp.Body.Close()
				
				if resp.StatusCode != http.StatusOK {
					results <- fmt.Errorf("unexpected status: %d", resp.StatusCode)
					continue
				}
				
				results <- nil
			}
		}()
	}

	for i := 0; i < concurrency*requests; i++ {
		if err := <-results; err != nil {
			t.Errorf("concurrent request failed: %v", err)
		}
	}
}

func loadTestBenchmark(t *testing.T, filename string) types.AnalyzeRequest {
	path := filepath.Join("test-ci", "mock-repos", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load test benchmark: %v", err)
	}

	var req types.AnalyzeRequest
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatalf("failed to parse benchmark data: %v", err)
	}

	return req
}

func loadWebhookPayload(t *testing.T, filename string) []byte {
	path := filepath.Join("test-ci", "github-sim", "webhook-payloads", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load webhook payload: %v", err)
	}
	return data
}

func loadBenchmarkScenario(t *testing.T, filename string) map[string]interface{} {
	path := filepath.Join("test-ci", "benchmark-data", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load benchmark scenario: %v", err)
	}

	var scenario map[string]interface{}
	if err := json.Unmarshal(data, &scenario); err != nil {
		t.Fatalf("failed to parse scenario data: %v", err)
	}

	return scenario
}

func validateScenario(t *testing.T, env *TestEnvironment, scenario map[string]interface{}) {
	t.Logf("Validating scenario: %s", scenario["scenario"])
}