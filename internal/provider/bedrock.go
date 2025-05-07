package provider

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	openai "github.com/sashabaranov/go-openai"
)

// BedrockProvider implements Provider for AWS Bedrock models.
// For simplicity, we treat modelId as the Provider name exposed via Config.model_name.
// This implementation supports text-generation models that accept JSON with "prompt" field
// and return JSON with "completion" or "generated_text".
// It is not production-ready but demonstrates wiring.

type BedrockProvider struct {
	name    string
	modelId string
	client  *bedrockruntime.Client
}

func NewBedrockProvider(name, modelId string) (*BedrockProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	brClient := bedrockruntime.NewFromConfig(cfg, func(o *bedrockruntime.Options) {
		// Set custom endpoint if needed via AWS_ENDPOINT_URL env
	})
	return &BedrockProvider{
		name:    name,
		modelId: modelId,
		client:  brClient,
	}, nil
}

func (p *BedrockProvider) Name() string { return p.name }

func messagesToPrompt(msgs []openai.ChatCompletionMessage) string {
	var b bytes.Buffer
	for _, m := range msgs {
		if m.Role == "user" {
			b.WriteString("User: ")
		} else if m.Role == "assistant" {
			b.WriteString("Assistant: ")
		} else {
			b.WriteString(m.Role + ": ")
		}
		b.WriteString(m.Content)
		b.WriteString("\n")
	}
	b.WriteString("Assistant: ")
	return b.String()
}

func (p *BedrockProvider) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	payload := map[string]interface{}{
		"prompt":      messagesToPrompt(req.Messages),
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"top_p":       req.TopP,
		"stop":        req.Stop,
	}
	bodyBytes, _ := json.Marshal(payload)

	out, err := p.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        bodyBytes,
		ModelId:     aws.String(p.modelId),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}
	respBytes := out.Body

	// Assume response JSON like {"completion": "text"} or {"generated_text": "text"}
	var respMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &respMap); err != nil {
		// fallback treat as plain string
		respMap = map[string]interface{}{"completion": string(respBytes)}
	}

	text := ""
	if v, ok := respMap["completion"].(string); ok {
		text = v
	} else if v, ok := respMap["generated_text"].(string); ok {
		text = v
	} else {
		text = string(respBytes)
	}

	choice := openai.ChatCompletionChoice{
		Index: 0,
		Message: openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: text,
		},
		FinishReason: "stop",
	}

	return &openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{choice},
		Usage:   openai.Usage{},
	}, nil
}

func (p *BedrockProvider) Embedding(ctx context.Context, req openai.EmbeddingRequest) (*openai.EmbeddingResponse, error) {
	return nil, &openai.APIError{Message: "embedding not implemented for bedrock"}
}
