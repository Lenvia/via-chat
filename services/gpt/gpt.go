package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	gogpt "github.com/sashabaranov/go-openai"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"net/url"
)

var OpenaiClient *openai.Client

func LoadGPT(file *ini.File) {
	ApiKey := file.Section("gpt").Key("API_KEY").String()

	config := gogpt.DefaultConfig(ApiKey)
	proxyUrl, err := url.Parse("http://127.0.0.1:7890") // WARNING 当前仅支持本地，如何在容器中穿透？
	// proxyUrl, err := url.Parse("http://host.docker.internal:7890")  // 访问宿主机代理
	if err != nil {
		log.Println(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	OpenaiClient = gogpt.NewClientWithConfig(config)
}

func GetReply(client *openai.Client, query string) string {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: query},
			},
		},
	)

	if err != nil {
		log.Println(err)
	}
	reply := resp.Choices[0].Message.Content
	return reply
}
