// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package regression

import (
	"fmt"

	"regression-ci/pkg/types"
)

func (d *Detector) getBaseline(repo, component string) (*types.Baseline, error) {
	var baseline types.Baseline
	query := `SELECT repo, component, baseline_value, sample_count, updated_at 
	          FROM baselines WHERE repo = ? AND component = ?`
	
	err := d.db.Get(&baseline, query, repo, component)
	if err != nil {
		return nil, fmt.Errorf("baseline not found: %w", err)
	}
	
	return &baseline, nil
}

func (d *Detector) getRecentSamples(repo, component string, limit int) ([]float64, error) {
	query := `SELECT value FROM benchmarks 
	          WHERE repo = ? AND component = ? 
	          ORDER BY timestamp DESC LIMIT ?`
	
	var values []float64
	err := d.db.Select(&values, query, repo, component, limit)
	return values, err
}

func (d *Detector) getRepoConfig(repo string) (*types.RepoConfig, error) {
	var config types.RepoConfig
	query := `SELECT repo, threshold_percent, min_samples, enabled 
	          FROM config WHERE repo = ?`
	
	err := d.db.Get(&config, query, repo)
	return &config, err
}