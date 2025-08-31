// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type WebhookSimulator struct {
	targetURL string
	secret    string
	client    *http.Client
}

func NewSimulator(targetURL, secret string) *WebhookSimulator {
	return &WebhookSimulator{
		targetURL: targetURL,
		secret:    secret,
		client:    &http.Client{},
	}
}

func (w *WebhookSimulator) SendPayload(payloadFile string) error {
	payload, err := w.loadPayload(payloadFile)
	if err != nil {
		return fmt.Errorf("failed to load payload: %w", err)
	}

	signature := w.generateSignature(payload)
	
	req, err := http.NewRequest("POST", w.targetURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "pull_request")
	req.Header.Set("X-Hub-Signature-256", signature)
	req.Header.Set("User-Agent", "GitHub-Hookshot/test")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %d %s\nBody: %s\n", resp.StatusCode, resp.Status, string(body))

	return nil
}

func (w *WebhookSimulator) loadPayload(filename string) ([]byte, error) {
	payloadPath := filepath.Join("test-ci", "github-sim", "webhook-payloads", filename)
	return os.ReadFile(payloadPath)
}

func (w *WebhookSimulator) generateSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(w.secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run simulator.go <target-url> <payload-file>")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	payloadFile := os.Args[2]
	secret := os.Getenv("WEBHOOK_SECRET")
	
	if secret == "" {
		secret = "test-secret-key"
		fmt.Println("Using default test secret")
	}

	simulator := NewSimulator(targetURL, secret)
	
	if err := simulator.SendPayload(payloadFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}