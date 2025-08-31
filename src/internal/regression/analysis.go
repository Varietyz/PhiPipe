// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package regression

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/stat"

	"regression-ci/pkg/types"
)

func (d *Detector) detectRegression(repo, component string, currentValue float64) (*types.RegressionResult, error) {
	baseline, err := d.getBaseline(repo, component)
	if err != nil {
		return d.createInitialBaseline(repo, component, currentValue)
	}

	percentChange := ((currentValue - baseline.BaselineValue) / baseline.BaselineValue) * 100
	
	threshold := d.config.DefaultThreshold
	repoConfig, err := d.getRepoConfig(repo)
	if err == nil && repoConfig.ThresholdPercent > 0 {
		threshold = repoConfig.ThresholdPercent
	}

	isRegression := percentChange > threshold
	confidence := d.calculateConfidence(baseline, currentValue, percentChange)

	return &types.RegressionResult{
		IsRegression:    isRegression,
		CurrentValue:    currentValue,
		BaselineValue:   baseline.BaselineValue,
		PercentChange:   percentChange,
		ConfidenceScore: confidence,
		SampleSize:      baseline.SampleCount,
	}, nil
}

func (d *Detector) createInitialBaseline(repo, component string, value float64) (*types.RegressionResult, error) {
	baseline := &types.Baseline{
		Repo:          repo,
		Component:     component,
		BaselineValue: value,
		SampleCount:   1,
		UpdatedAt:     time.Now().Unix(),
	}

	query := `INSERT OR REPLACE INTO baselines 
	          (repo, component, baseline_value, sample_count, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`
	
	_, err := d.db.Exec(query, baseline.Repo, baseline.Component, 
		baseline.BaselineValue, baseline.SampleCount, baseline.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create baseline: %w", err)
	}

	return &types.RegressionResult{
		IsRegression:    false,
		CurrentValue:    value,
		BaselineValue:   value,
		PercentChange:   0.0,
		ConfidenceScore: 100.0,
		SampleSize:      1,
	}, nil
}

func (d *Detector) updateBaseline(repo, component string, newValue float64, result *types.RegressionResult) {
	if result.IsRegression {
		return
	}

	recentSamples, err := d.getRecentSamples(repo, component, d.config.MaxSamples)
	if err != nil || len(recentSamples) < d.config.MinSamples {
		return
	}

	newBaseline := stat.Mean(recentSamples, nil)
	
	query := `UPDATE baselines SET baseline_value = ?, sample_count = ?, updated_at = ? 
	          WHERE repo = ? AND component = ?`
	
	d.db.Exec(query, newBaseline, len(recentSamples), time.Now().Unix(), repo, component)
}

func (d *Detector) calculateConfidence(baseline *types.Baseline, currentValue, percentChange float64) float64 {
	if baseline.SampleCount < d.config.MinSamples {
		return 50.0
	}

	changeAbs := math.Abs(percentChange)
	confidence := math.Min(90.0, 50.0+(changeAbs*2))
	
	return confidence
}