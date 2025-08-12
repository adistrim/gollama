package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"gollama/config"

	"github.com/google/go-github/v74/github"
	"github.com/sashabaranov/go-openai"
)

func createGitHubBranchTool() Tool {
	return Tool{
		Definition: openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "create_github_branch",
				Description: "Create a new branch in a GitHub repository.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"owner": map[string]any{
							"type":        "string",
							"description": "The owner or organization of the repository.",
						},
						"repo": map[string]any{
							"type":        "string",
							"description": "The name of the repository.",
						},
						"branch_name": map[string]any{
							"type":        "string",
							"description": "The name for the new branch.",
						},
						"source_branch": map[string]any{
							"type":        "string",
							"description": "The source branch to create from (defaults to 'main').",
						},
					},
					"required": []string{"owner", "repo", "branch_name"},
				},
			},
		},
		Execute: func(ctx context.Context, args string) (string, error) {
			type branchArgs struct {
				Owner        string `json:"owner"`
				Repo         string `json:"repo"`
				BranchName   string `json:"branch_name"`
				SourceBranch string `json:"source_branch"`
			}

			var parsedArgs branchArgs
			err := json.Unmarshal([]byte(args), &parsedArgs)
			if err != nil {
				return "", fmt.Errorf("failed to parse tool arguments: %w", err)
			}

			if parsedArgs.SourceBranch == "" {
				parsedArgs.SourceBranch = "main"
			}

			client := github.NewClient(nil).WithAuthToken(config.ENV.GithubToken)

			sourceRef, _, err := client.Git.GetRef(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				fmt.Sprintf("heads/%s", parsedArgs.SourceBranch),
			)
			if err != nil {
				return "", fmt.Errorf("failed to get source branch reference: %w", err)
			}

			newRef := &github.Reference{
			    Ref: github.Ptr(fmt.Sprintf("refs/heads/%s", parsedArgs.BranchName)),
			    Object: &github.GitObject{
			        SHA: sourceRef.Object.SHA,
			    },
			}

			createdRef, _, err := client.Git.CreateRef(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				newRef,
			)
			if err != nil {
				return "", fmt.Errorf("failed to create branch: %w", err)
			}

			result := map[string]any{
				"branch_name": parsedArgs.BranchName,
				"sha":         createdRef.Object.GetSHA(),
				"ref":         createdRef.GetRef(),
				"url":         createdRef.GetURL(),
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal branch result: %w", err)
			}

			return string(resultBytes), nil
		},
	}
}
