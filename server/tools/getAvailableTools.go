package tools

import (
	"context"
	
	"github.com/sashabaranov/go-openai"
)

type Tool struct {
	Definition openai.Tool
	Execute func(ctx context.Context, args string) (string, error)
}

func GetAvailableTools() map[string]Tool {
    tools := make(map[string]Tool)
    tools["get_github_issue_details"] = getGitHubIssueDetailsTool()
    tools["create_github_pr"] = createGitHubPRTool()
    tools["create_github_branch"] = createGitHubBranchTool()
    tools["get_repository_files"] = getRepositoryFilesTool()
    tools["update_github_file"] = updateGitHubFileTool()
    return tools
}
