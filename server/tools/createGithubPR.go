package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"gollama/config"

	"github.com/google/go-github/v74/github"
	"github.com/sashabaranov/go-openai"
)

func createGitHubPRTool() Tool {
	return Tool{
		Definition: openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "create_github_pr",
				Description: "Create a new pull request in a GitHub repository.",
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
						"title": map[string]any{
							"type":        "string",
							"description": "The title of the pull request.",
						},
						"body": map[string]any{
							"type":        "string",
							"description": "The body/description of the pull request.",
						},
						"head": map[string]any{
							"type":        "string",
							"description": "The name of the branch where your changes are implemented (source branch).",
						},
						"base": map[string]any{
							"type":        "string",
							"description": "The name of the branch you want the changes pulled into (target branch, usually 'main' or 'master').",
						},
						"draft": map[string]any{
							"type":        "boolean",
							"description": "Whether to create the PR as a draft. Defaults to false.",
						},
					},
					"required": []string{"owner", "repo", "title", "body", "head", "base"},
				},
			},
		},
		Execute: func(ctx context.Context, args string) (string, error) {
			type prArgs struct {
				Owner string `json:"owner"`
				Repo  string `json:"repo"`
				Title string `json:"title"`
				Body  string `json:"body"`
				Head  string `json:"head"`
				Base  string `json:"base"`
				Draft bool   `json:"draft"`
			}

			var parsedArgs prArgs
			err := json.Unmarshal([]byte(args), &parsedArgs)
			if err != nil {
				return "", fmt.Errorf("failed to parse tool arguments: %w", err)
			}

			client := github.NewClient(nil).WithAuthToken(config.ENV.GithubToken)

			newPR := &github.NewPullRequest{
				Title: &parsedArgs.Title,
				Head:  &parsedArgs.Head,
				Base:  &parsedArgs.Base,
				Body:  &parsedArgs.Body,
				Draft: &parsedArgs.Draft,
			}

			pr, _, err := client.PullRequests.Create(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				newPR,
			)
			if err != nil {
				return "", fmt.Errorf("failed to create pull request via GitHub API: %w", err)
			}

			result := map[string]any{
				"number":     pr.GetNumber(),
				"title":      pr.GetTitle(),
				"state":      pr.GetState(),
				"url":        pr.GetHTMLURL(),
				"author":     pr.GetUser().GetLogin(),
				"head":       pr.GetHead().GetRef(),
				"base":       pr.GetBase().GetRef(),
				"draft":      pr.GetDraft(),
				"created_at": pr.GetCreatedAt(),
				"body":       pr.GetBody(),
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal PR result: %w", err)
			}

			return string(resultBytes), nil
		},
	}
}
