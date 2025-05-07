package router

import (
	"context"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

type stubProvider struct{ id string }

func (s stubProvider) Name() string { return s.id }
func (s stubProvider) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	return nil, nil
}
func (s stubProvider) Embedding(ctx context.Context, req openai.EmbeddingRequest) (*openai.EmbeddingResponse, error) {
	return nil, nil
}

func TestRouterRoundRobin(t *testing.T) {
	r := New()
	p1 := stubProvider{"p1"}
	p2 := stubProvider{"p2"}

	r.Register("model", p1)
	r.Register("model", p2)

	for i := 0; i < 10; i++ {
		p, ok := r.Get("model")
		if !ok {
			t.Fatalf("router returned false")
		}
		expected := "p1"
		if i%2 == 1 {
			expected = "p2"
		}
		if p.Name() != expected {
			t.Fatalf("round robin failed at %d expected %s got %s", i, expected, p.Name())
		}
	}
}
