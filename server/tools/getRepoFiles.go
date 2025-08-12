package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"gollama/config"

	"github.com/google/go-github/v74/github"
	"github.com/sashabaranov/go-openai"
)

func getRepositoryFilesTool() Tool {
	return Tool{
		Definition: openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_repository_files",
				Description: "Get repository structure and file contents to understand the codebase.",
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
							"description": "The path to explore (empty for root, or specific directory/file path).",
						},
						"ref": map[string]any{
							"type":        "string",
							"description": "The branch/commit reference (defaults to default branch).",
						},
					},
					"required": []string{"owner", "repo"},
				},
			},
		},
		Execute: func(ctx context.Context, args string) (string, error) {
			type repoArgs struct {
				Owner string `json:"owner"`
				Repo  string `json:"repo"`
				Path  string `json:"path"`
				Ref   string `json:"ref"`
			}

			var parsedArgs repoArgs
			err := json.Unmarshal([]byte(args), &parsedArgs)
			if err != nil {
				return "", fmt.Errorf("failed to parse tool arguments: %w", err)
			}

			client := github.NewClient(nil).WithAuthToken(config.ENV.GithubToken)

			opts := &github.RepositoryContentGetOptions{}
			if parsedArgs.Ref != "" {
				opts.Ref = parsedArgs.Ref
			}

			fileContent, directoryContent, _, err := client.Repositories.GetContents(
				ctx,
				parsedArgs.Owner,
				parsedArgs.Repo,
				parsedArgs.Path,
				opts,
			)
			if err != nil {
				return "", fmt.Errorf("failed to get repository contents: %w", err)
			}

			result := map[string]any{}

			if fileContent != nil {
				content, err := fileContent.GetContent()
				if err != nil {
					return "", fmt.Errorf("failed to decode file content: %w", err)
				}
				result = map[string]any{
					"type":     "file",
					"name":     fileContent.GetName(),
					"path":     fileContent.GetPath(),
					"content":  content,
					"size":     fileContent.GetSize(),
					"sha":      fileContent.GetSHA(),
				}
			} else {
				var files []map[string]any
				for _, item := range directoryContent {
					files = append(files, map[string]any{
						"name": item.GetName(),
						"path": item.GetPath(),
						"type": item.GetType(),
						"size": item.GetSize(),
						"sha":  item.GetSHA(),
					})
				}
				result = map[string]any{
					"type":  "directory",
					"path":  parsedArgs.Path,
					"files": files,
				}
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal repository result: %w", err)
			}

			return string(resultBytes), nil
		},
	}
}
