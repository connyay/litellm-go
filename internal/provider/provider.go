package provider

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

// Provider wraps a backing LLM endpoint (OpenAI, Azure, Bedrock, etc.).
// It takes an OpenAI-style ChatCompletionRequest and returns an OpenAI-style response.
// Streaming not yet supported.
//
// Implementations must be goroutine-safe.
//
// The abstraction lets the router treat every upstream uniformly.
// Future: add Embeddings, Images, etc.

type Provider interface {
	ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
	Name() string // for logging/debug
}
