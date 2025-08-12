package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"gollama/config"

	"github.com/google/go-github/v74/github"
	"github.com/sashabaranov/go-openai"
)

func updateGitHubFileTool() Tool {
	return Tool{
		Definition: openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "update_github_file",
				Description: "Create or update a file in a GitHub repository.",
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
						"path": map[string]any{
							"type":        "string",
							"description": "The file path in the repository.",
						},
						"content": map[string]any{
						    "type":        "string",
						    "description": "The full, updated file content (not just the diff). Always include the entire file.",
						},
						"message": map[string]any{
							"type":        "string",
							"description": "The commit message for this change.",
						},
						"branch": map[string]any{
							"type":        "string",
							"description": "The branch to commit to.",
						},
					},
					"required": []string{"owner", "repo", "path", "content", "message", "branch"},
				},
			},
		},
		Execute: func(ctx context.Context, args string) (string, error) {
			type fileArgs struct {
				Owner   string `json:"owner"`
				Repo    string `json:"repo"`
				Path    string `json:"path"`
				Content string `json:"content"`
				Message string `json:"message"`
				Branch  string `json:"branch"`
			}

			var parsedArgs fileArgs
			err := json.Unmarshal([]byte(args), &parsedArgs)
			if err != nil {
				return "", fmt.Errorf("failed to parse tool arguments: %w", err)
			}

			client := github.NewClient(nil).WithAuthToken(config.ENV.GithubToken)

			opts := &github.RepositoryContentFileOptions{
				Message: &parsedArgs.Message,
				Content: []byte(parsedArgs.Content),
				Branch:  &parsedArgs.Branch,
			}

			existingFile, _, _, err := client.Repositories.GetContents(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				parsedArgs.Path,
				&github.RepositoryContentGetOptions{Ref: parsedArgs.Branch},
			)
			if err == nil && existingFile != nil {
				opts.SHA = existingFile.SHA
			}

			fileResponse, _, err := client.Repositories.CreateFile(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				parsedArgs.Path,
				opts,
			)
			if err != nil {
				return "", fmt.Errorf("failed to update file: %w", err)
			}

			result := map[string]any{
				"path":       parsedArgs.Path,
				"sha":        fileResponse.Content.GetSHA(),
				"commit_sha": fileResponse.Commit.GetSHA(),
				"message":    parsedArgs.Message,
				"branch":     parsedArgs.Branch,
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal file update result: %w", err)
			}

			return string(resultBytes), nil
		},
	}
}
