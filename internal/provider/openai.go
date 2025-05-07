package provider

import (
	"context"
	"net/http"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements Provider backed by OpenAI-compatible endpoint
// This also supports Azure OpenAI by passing an apiVersion and Azure deployment name.

type OpenAIProvider struct {
	name   string
	client *openai.Client
}

// NewOpenAIProvider constructs provider. If apiKey is empty, client uses env vars.
func NewOpenAIProvider(name, apiBase, apiKey, apiVersion string, azure bool) *OpenAIProvider {
	cfg := openai.DefaultConfig(apiKey)
	if apiBase != "" {
		cfg.BaseURL = apiBase
	}
	// For Azure, BaseURL should include /openai before /deployments, but we expect user to set correctly.
	if azure {
		cfg.APIVersion = apiVersion
		cfg.APIType = openai.APITypeAzure
	}
	// Increase timeout a bit
	cfg.HTTPClient = &http.Client{Timeout: 60 * time.Second}

	return &OpenAIProvider{
		name:   name,
		client: openai.NewClientWithConfig(cfg),
	}
}

func (p *OpenAIProvider) Name() string { return p.name }

func (p *OpenAIProvider) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
