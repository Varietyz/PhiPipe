// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package regression

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"regression-ci/internal/config"
	"regression-ci/pkg/types"
)

type Detector struct {
	db     *sqlx.DB
	config config.DetectionConfig
}

func New(db *sqlx.DB, cfg config.DetectionConfig) *Detector {
	return &Detector{
		db:     db,
		config: cfg,
	}
}

func (d *Detector) Analyze(req types.AnalyzeRequest) (*types.AnalyzeResponse, error) {
	timestamp := time.Now().Unix()
	response := &types.AnalyzeResponse{
		Repo:      req.Repo,
		Commit:    req.Commit,
		Timestamp: timestamp,
	}

	if err := d.storeBenchmarks(req, timestamp); err != nil {
		return nil, fmt.Errorf("failed to store benchmarks: %w", err)
	}

	for component, value := range req.Components {
		result, err := d.detectRegression(req.Repo, component, value)
		
		componentResult := types.ComponentResult{
			Component: component,
		}
		
		if err != nil {
			componentResult.Error = err.Error()
		} else {
			componentResult.Result = result
			d.updateBaseline(req.Repo, component, value, result)
		}
		
		response.Components = append(response.Components, componentResult)
	}

	return response, nil
}

func (d *Detector) storeBenchmarks(req types.AnalyzeRequest, timestamp int64) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO benchmarks (repo, branch, commit_hash, component, value, timestamp) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	for component, value := range req.Components {
		_, err := tx.Exec(query, req.Repo, req.Branch, req.Commit, component, value, timestamp)
		if err != nil {
			return fmt.Errorf("failed to insert benchmark: %w", err)
		}
	}

	return tx.Commit()
}