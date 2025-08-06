package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"gollama/config"

	"github.com/google/go-github/v74/github"
	"github.com/sashabaranov/go-openai"
)

type Tool struct {
	Definition openai.Tool
	Execute    func(args string) (string, error)
}

func GetAvailableTools() map[string]Tool {
	tools := make(map[string]Tool)
	tools["get_github_issue_details"] = getGitHubIssueDetailsTool()
	return tools
}

func getGitHubIssueDetailsTool() Tool {
	return Tool{
		Definition: openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_github_issue_details",
				Description: "Get detailed information about a specific issue from a GitHub repository.",
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
				        "issue_number": map[string]any{
				            "type":        "integer",
				            "description": "The number of the issue to retrieve.",
				        },
				    },
				    "required": []string{"owner", "repo", "issue_number"},
				},
			},
		},
		Execute: func(args string) (string, error) {
			type issueArgs struct {
				Owner       string `json:"owner"`
				Repo        string `json:"repo"`
				IssueNumber json.Number `json:"issue_number"`
			}

			var parsedArgs issueArgs
			err := json.Unmarshal([]byte(args), &parsedArgs)
			if err != nil {
				return "", fmt.Errorf("failed to parse tool arguments: %w", err)
			}

			issueNum, err := parsedArgs.IssueNumber.Int64()
			if err != nil {
				return "", fmt.Errorf("failed to parse issue number: %w", err)
			}
			
			client := github.NewClient(nil).WithAuthToken(config.ENV.GithubToken)
			
			issue, _, err := client.Issues.Get(
				context.Background(),
				parsedArgs.Owner,
				parsedArgs.Repo,
				int(issueNum),
			)
			if err != nil {
				return "", fmt.Errorf("failed to get issue from GitHub API: %w", err)
			}
			
			result := map[string]any{
				"title": issue.GetTitle(),
				"state": issue.GetState(),
				"author": issue.GetUser().GetLogin(),
				"body": issue.GetBody(),
				"labels": issue.Labels,
				"url": issue.GetHTMLURL(),
				"created_at": issue.GetCreatedAt(),
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal issue result: %w", err)
			}

			return string(resultBytes), nil
		},
	}
}
