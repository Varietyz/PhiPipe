// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"

	"regression-ci/internal/config"
)

type Client struct {
	client *github.Client
	config config.GitHubConfig
}

func New(cfg config.GitHubConfig) *Client {
	var client *github.Client
	
	if cfg.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &Client{
		client: client,
		config: cfg,
	}
}

func (c *Client) CreatePRComment(ctx context.Context, owner, repo string, prNumber int, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := c.client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to create PR comment: %w", err)
	}

	return nil
}

func (c *Client) UpdatePRComment(ctx context.Context, owner, repo string, commentID int64, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := c.client.Issues.EditComment(ctx, owner, repo, commentID, comment)
	if err != nil {
		return fmt.Errorf("failed to update PR comment: %w", err)
	}

	return nil
}

func (c *Client) ListPRComments(ctx context.Context, owner, repo string, prNumber int) ([]*github.IssueComment, error) {
	comments, _, err := c.client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list PR comments: %w", err)
	}

	return comments, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return repository, nil
}

func (c *Client) ValidateWebhookSignature(payload []byte, signature string) bool {
	if c.config.WebhookSecret == "" {
		return false
	}

	err := github.ValidateSignature(signature, payload, []byte(c.config.WebhookSecret))
	return err == nil
}