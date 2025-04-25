package core

import (
	"context"
	"github.com/leslieleung/ptpt/internal/config"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (string, Usage, error)
	StreamChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error)
	SetTemperature(t float32)
}

func GetClient() Client {
	cfg := config.GetIns()
	ai := cfg.AiName
	log.Debugf("当前使用ai: %v", ai)
	switch ai {
	case "kimi":
		return &Kimi{}
	case "qwen":
		return &TongYi{}
	case "deepseek":
		return &DeepSeek{}
	default:
		return &OpenAI{}
	}
}
