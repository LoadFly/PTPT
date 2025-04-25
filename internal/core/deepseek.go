package core

import (
	"context"
	"github.com/avast/retry-go/v4"
	"github.com/leslieleung/ptpt/internal/config"
	"github.com/leslieleung/ptpt/internal/interract"
	"github.com/leslieleung/ptpt/internal/ui"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
)

type DeepSeek struct {
	client      *openai.Client
	once        sync.Once
	temperature float32
}

var deepseekDefaultUrl = "https://api.deepseek.com/v1"

func (d *DeepSeek) getClient() *openai.Client {
	cfg := config.GetIns()
	if cfg.APIKey == "" {
		ui.ErrorfExit("API key is not set. Please set it in %s", filepath.Join(interract.GetPTPTDir(), "config.yaml"))
	}
	if Model == "" {
		Model = "deepseek-chat"  // 默认模型，需要根据实际的DeepSeek API调整
	}
	d.once.Do(func() {
		c := openai.DefaultConfig(cfg.APIKey)
		c.BaseURL = deepseekDefaultUrl
		if cfg.Proxy != "" {
			proxy, _ := url.Parse(cfg.Proxy)
			c.HTTPClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
		d.client = openai.NewClientWithConfig(c)
	})
	d.temperature = Temperature
	return d.client
}

func (d *DeepSeek) CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (string, Usage, error) {
	var resp openai.ChatCompletionResponse
	var err error
	err = retry.Do(func() error {
		resp, err = d.getClient().CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       Model,
			Messages:    messages,
			Temperature: d.temperature,
		})
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(1))
	log.Debugf("DeepSeek Token Usage [Prompt: %d, Completion: %d, Total: %d]",
		resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	log.Debugf("Response: %+v", resp)
	if len(resp.Choices) == 0 {
		return "", Usage{}, nil
	}
	return resp.Choices[0].Message.Content, Usage(resp.Usage), nil
}

func (d *DeepSeek) StreamChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	return d.getClient().CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    Model,
		Messages: messages,
		Stream:   true,
	})
}

func (d *DeepSeek) SetTemperature(t float32) {
	d.temperature = t
}