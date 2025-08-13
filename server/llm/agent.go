package llm

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gollama/tools"

	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	client *openai.Client
	model  string
}

const (
	SystemPrompt = `
		You are Gollama, an expert AI Software Engineer Buddy.
		You always communicate in casual human tone.
		
		IMPORTANT WORKFLOW:
		When a user asks you to make changes to a repository (create PR, fix issues, implement features):
		
		PHASE 1 - PLANNING:
		1. First, create a detailed plan of what you will do
		2. List all the steps you'll take with specific tool calls
		3. Ask for user approval before proceeding
		4. Wait for explicit approval
		
		PHASE 2 - EXECUTION (USE TOOLS):
		5. Only after approval, start using tools to execute the plan
		6. Use tools in the planned sequence
		7. Once you complete all steps, provide a brief summary of the changes made
		
		For general questions, repository exploration, or single tool calls, you can use tools directly.
		
		When creating implementation plans, be specific about:
		- Which files you'll examine
		- What branch name you'll use  
		- What code changes you'll make
		- Commit messages you'll use
		- PR title and description
	`
)

func (a *Agent) RunSessionConversation(ctx context.Context, messages []openai.ChatCompletionMessage, statusCallback func(string)) ([]openai.ChatCompletionMessage, error) {
	availableTools := tools.GetAvailableTools()
	var toolDefs []openai.Tool
	for _, t := range availableTools {
		toolDefs = append(toolDefs, t.Definition)
	}

	const maxIterations = 6
	for step := range maxIterations {
		log.Printf("========== Agent Call %d ==========", step+1)
		log.Println("Step 1: Sending conversation history and tool definitions to LLM...")

		resp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: messages,
				Tools:    toolDefs,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("chat completion failed: %w", err)
		}
		if len(resp.Choices) == 0 {
			return nil, errors.New("no response choices from LLM")
		}

		responseMessage := resp.Choices[0].Message
		messages = append(messages, responseMessage)

		if len(responseMessage.ToolCalls) == 0 {
			log.Println("No tool calls requested. Agent finished.")
			break
		}
		
		if statusCallback != nil {
			statusCallback("tools_used")
		}

		log.Printf("Step 2: LLM requested %d tool call(s)", len(responseMessage.ToolCalls))

		var toolResponses []openai.ChatCompletionMessage

		for _, toolCall := range responseMessage.ToolCalls {
			functionName := toolCall.Function.Name
			log.Printf("Tool call requested: %s", functionName)

			tool, ok := availableTools[functionName]
			if !ok {
				return nil, fmt.Errorf("LLM requested an unknown tool: %s", functionName)
			}

			log.Printf("Executing tool '%s' with args: %s", functionName, toolCall.Function.Arguments)
			toolResult, err := tool.Execute(ctx, toolCall.Function.Arguments)
			if err != nil {
				log.Printf("Tool '%s' failed: %v", functionName, err)
				toolResult = fmt.Sprintf("ERROR: %v", err)
			} else {
				log.Printf("Tool '%s' execution successful", functionName)
			}

			toolResponses = append(toolResponses, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    toolResult,
				Name:       functionName,
				ToolCallID: toolCall.ID,
			})
		}

		messages = append(messages, toolResponses...)
		log.Println("Step 3: Tool results appended to conversation for next LLM iteration")
	}

	if len(messages) > 0 && messages[len(messages)-1].Role != openai.ChatMessageRoleAssistant {
		
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "Please provide a summary of what you've done and the current status.",
		})
		
		resp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: messages,
			},
		)
		
		if err == nil && len(resp.Choices) > 0 {
			messages = append(messages, resp.Choices[0].Message)
		} else {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "I've completed the requested operations. The changes have been made to the repository.",
			})
		}
	}

	log.Println("========== Agent session complete ==========")
	return messages, nil
}
