package core

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/leslieleung/ptpt/internal/config"
	"github.com/leslieleung/ptpt/internal/interract"
	"github.com/leslieleung/ptpt/internal/ui"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

// aliyun dashscope support
type TongYi struct {
	client      *openai.Client
	once        sync.Once
	temperature float32
}

const TongYiDefaultUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1"

func (k *TongYi) getClient() *openai.Client {
	cfg := config.GetIns()
	if cfg.APIKey == "" {
		ui.ErrorfExit("API key is not set. Please set it in %s", filepath.Join(interract.GetPTPTDir(), "config.yaml"))
	}
	if Model == "" {
		Model = "qwen-turbo"
	}
	k.once.Do(func() {
		c := openai.DefaultConfig(cfg.APIKey)
		c.BaseURL = TongYiDefaultUrl
		if cfg.Proxy != "" {
			proxy, _ := url.Parse(cfg.Proxy)
			c.HTTPClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
		k.client = openai.NewClientWithConfig(c)
	})
	k.temperature = Temperature
	return k.client
}

func (k *TongYi) CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (string, Usage, error) {
	var resp openai.ChatCompletionResponse
	var err error
	err = retry.Do(func() error {
		resp, err = k.getClient().CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       Model,
			Messages:    messages,
			Temperature: k.temperature,
		})
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(1))
	log.Debugf("TongYi Token Usage [Prompt: %d, Completion: %d, Total: %d]",
		resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	log.Debugf("Response: %+v", resp)
	if len(resp.Choices) == 0 {
		return "", Usage{}, nil
	}
	return resp.Choices[0].Message.Content, Usage(resp.Usage), nil
}

func (k *TongYi) StreamChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	return k.getClient().CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    Model,
		Messages: messages,
		Stream:   true,
	})
}

func (k *TongYi) SetTemperature(t float32) {
	k.temperature = t
}
