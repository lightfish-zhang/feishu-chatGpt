package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	BASEURL     = "https://api.openai.com/v1/"
	maxTokens   = 2048
	temperature = 0.3
	top_p       = 1
	engine      = "text-davinci-003"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

func Completions(msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            engine,
		Prompt:           msg,
		MaxTokens:        maxTokens,
		Temperature:      temperature,
		TopP:             top_p,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := viper.GetString("OPENAI_KEY")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		return "", fmt.Errorf("gtp api %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Text
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}

func FormatQuestion(question string) string {
	return "Answer:" + question
}
