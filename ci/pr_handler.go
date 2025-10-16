package ci

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// CommentOnPR posts a comment on the current PR using GitHub CLI (`gh`).
func CommentOnPR(comment string) error {
	prNumber := os.Getenv("GITHUB_PR_NUMBER")
	repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("GITHUB_TOKEN")

	if prNumber == "" || repo == "" || token == "" {
		return fmt.Errorf("GITHUB_PR_NUMBER, GITHUB_REPOSITORY or GITHUB_TOKEN not set")
	}

	// Prepare GitHub CLI command
	cmd := exec.Command("gh", "pr", "comment", prNumber, "--repo", repo, "--body", comment)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GH_TOKEN=%s", token))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to post PR comment: %v, stderr: %s", err, stderr.String())
	}

	fmt.Printf("PR comment posted with response: %s\n", out.String())
	return nil
}