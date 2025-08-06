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

func (a *Agent) RunConversation(ctx context.Context, userInput string) (string, error) {
	availableTools := tools.GetAvailableTools()
	var toolDefs []openai.Tool
	for _, t := range availableTools {
		toolDefs = append(toolDefs, t.Definition)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: "Your name is Gollama, you help with tasks related to GitHub and provide information about GitHub and its features. If a question is completely unrelated to GitHub, say 'I'm sorry, I can't help with that.' Always respond in plain text - no markdown.",
		},
		{
			Role: openai.ChatMessageRoleUser,
			Content: userInput,
		},
	}

	log.Println("Step 1: Sending user prompt and tool definitions to LLM...")
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    a.model,
			Messages: messages,
			Tools:    toolDefs,
		},
	)

	if err != nil {
		return "", fmt.Errorf("initial chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices from LLM")
	}

	responseMessage := resp.Choices[0].Message
	messages = append(messages, responseMessage)

	if len(responseMessage.ToolCalls) > 0 {
		log.Println("Step 2: LLM wants to call a tool.")
		toolCall := responseMessage.ToolCalls[0]
		functionName := toolCall.Function.Name
		
		log.Printf("Tool call requested: %s", functionName)

		tool, ok := availableTools[functionName]
		if !ok {
			return "", fmt.Errorf("LLM requested an unknown tool: %s", functionName)
		}

		log.Printf("Executing tool '%s' with args: %s", functionName, toolCall.Function.Arguments)
		toolResult, err := tool.Execute(ctx, toolCall.Function.Arguments)
		if err != nil {
			return "", fmt.Errorf("failed to execute tool %s: %w", functionName, err)
		}

		log.Println("Step 3: Sending tool result back to LLM for final response.")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResult,
			Name:       functionName,
			ToolCallID: toolCall.ID,
		})

		finalResp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: messages,
			},
		)
		if err != nil {
			return "", fmt.Errorf("final chat completion failed: %w", err)
		}

		if len(finalResp.Choices) == 0 {
			return "", errors.New("no final response choices from LLM")
		}

		return finalResp.Choices[0].Message.Content, nil
	}
	
	log.Println("Step 2: LLM responded directly without a tool call.")
	return responseMessage.Content, nil
}
