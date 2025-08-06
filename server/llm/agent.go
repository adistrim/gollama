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
		You are Gollama, an expert AI Software Engineer.
		You always communicate in casual human tone - no emojis.
		
		IMPORTANT: Only use tools when the user explicitly provides ALL required information:
		- Do NOT use tools for general questions, introductions, or when missing required parameters
		- If the user asks general questions about your capabilities, respond directly without using any tools
	`
)

func (a *Agent) RunSessionConversation(ctx context.Context, messages []openai.ChatCompletionMessage) ([]openai.ChatCompletionMessage, error) {
	availableTools := tools.GetAvailableTools()
	var toolDefs []openai.Tool
	for _, t := range availableTools {
		toolDefs = append(toolDefs, t.Definition)
	}

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
		return nil, fmt.Errorf("initial chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response choices from LLM")
	}

	responseMessage := resp.Choices[0].Message
	updatedMessages := append(messages, responseMessage)

	if len(responseMessage.ToolCalls) > 0 {
		log.Println("Step 2: LLM wants to call a tool.")
		toolCall := responseMessage.ToolCalls[0]
		functionName := toolCall.Function.Name
		
		log.Printf("Tool call requested: %s", functionName)

		tool, ok := availableTools[functionName]
		if !ok {
			return nil, fmt.Errorf("LLM requested an unknown tool: %s", functionName)
		}

		log.Printf("Executing tool '%s' with args: %s", functionName, toolCall.Function.Arguments)
		toolResult, err := tool.Execute(ctx, toolCall.Function.Arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to execute tool %s: %w", functionName, err)
		}

		log.Println("Step 3: Sending tool result back to LLM for final response.")
		updatedMessages = append(updatedMessages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResult,
			Name:       functionName,
			ToolCallID: toolCall.ID,
		})

		finalResp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: updatedMessages,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("final chat completion failed: %w", err)
		}

		if len(finalResp.Choices) == 0 {
			return nil, errors.New("no final response choices from LLM")
		}

		finalMessage := finalResp.Choices[0].Message
		updatedMessages = append(updatedMessages, finalMessage)
	}
	
	return updatedMessages, nil
}
