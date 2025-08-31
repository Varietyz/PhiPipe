// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package types

type AnalyzeRequest struct {
	Repo       string                 `json:"repo" binding:"required"`
	Branch     string                 `json:"branch" binding:"required"`
	Commit     string                 `json:"commit" binding:"required"`
	Components map[string]float64     `json:"components" binding:"required"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type RegressionResult struct {
	IsRegression    bool    `json:"is_regression"`
	CurrentValue    float64 `json:"current_value"`
	BaselineValue   float64 `json:"baseline_value"`
	PercentChange   float64 `json:"percent_change"`
	ConfidenceScore float64 `json:"confidence_score"`
	SampleSize      int     `json:"sample_size"`
}

type ComponentResult struct {
	Component string            `json:"component"`
	Result    *RegressionResult `json:"result"`
	Error     string            `json:"error,omitempty"`
}

type AnalyzeResponse struct {
	Repo       string            `json:"repo"`
	Commit     string            `json:"commit"`
	Components []ComponentResult `json:"components"`
	Timestamp  int64             `json:"timestamp"`
}

type Baseline struct {
	Repo          string  `json:"repo" db:"repo"`
	Component     string  `json:"component" db:"component"`
	BaselineValue float64 `json:"baseline_value" db:"baseline_value"`
	SampleCount   int     `json:"sample_count" db:"sample_count"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
}

type Benchmark struct {
	ID         int64   `json:"id" db:"id"`
	Repo       string  `json:"repo" db:"repo"`
	Branch     string  `json:"branch" db:"branch"`
	CommitHash string  `json:"commit_hash" db:"commit_hash"`
	Component  string  `json:"component" db:"component"`
	Value      float64 `json:"value" db:"value"`
	Timestamp  int64   `json:"timestamp" db:"timestamp"`
}

type RepoConfig struct {
	Repo             string                        `json:"repo" db:"repo"`
	ThresholdPercent float64                       `json:"threshold_percent" db:"threshold_percent"`
	MinSamples       int                           `json:"min_samples" db:"min_samples"`
	Enabled          bool                          `json:"enabled" db:"enabled"`
	Components       map[string]ComponentConfig   `json:"components,omitempty"`
}

type ComponentConfig struct {
	CustomThreshold *float64 `json:"custom_threshold,omitempty"`
	Enabled         bool     `json:"enabled"`
}