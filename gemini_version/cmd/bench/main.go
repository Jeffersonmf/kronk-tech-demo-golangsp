package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/kronk/sdk/kronk/model"

	"github.com/jeff/gemini-kronk-demo/internal/kronkhub"
)

const prompt = "Explain in one paragraph how to fix a Go race condition using sync.Mutex."

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	path := kronkhub.ChatModelGPTOSS
	name := "gpt-oss-20b-Q8_0"

	if os.Getenv("BENCH_MODEL") == "qwen3" {
		path = kronkhub.ChatModelQwen3
		name = "Qwen3-8B-Q8_0"
	}

	fmt.Printf("=== %s ===\n", name)

	start := time.Now()
	krn, err := kronkhub.NewModel(path, model.WithContextWindow(8192))
	if err != nil {
		return err
	}
	loadDur := time.Since(start)
	fmt.Printf("load time: %s\n", loadDur)

	defer func() {
		if err := krn.Unload(context.Background()); err != nil {
			fmt.Printf("unload error: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	d := model.D{
		"messages":   []model.D{model.TextMessage("user", prompt)},
		"max_tokens": 512,
	}

	chatStart := time.Now()
	ch, err := krn.ChatStreaming(ctx, d)
	if err != nil {
		return fmt.Errorf("chat streaming: %w", err)
	}

	var last model.ChatResponse
	for resp := range ch {
		last = resp
		if resp.Choices[0].FinishReason() == model.FinishReasonError {
			return fmt.Errorf("error from model: %s", resp.Choices[0].Delta.Content)
		}
	}
	chatDur := time.Since(chatStart)

	if last.Usage == nil {
		return fmt.Errorf("no usage info in final response")
	}

	u := last.Usage
	fmt.Printf("chat wall time   : %s\n", chatDur)
	fmt.Printf("prompt tokens    : %d\n", u.PromptTokens)
	fmt.Printf("reasoning tokens : %d\n", u.ReasoningTokens)
	fmt.Printf("completion tokens: %d\n", u.CompletionTokens)
	fmt.Printf("output tokens    : %d\n", u.OutputTokens)
	fmt.Printf("tokens per second: %.2f\n", u.TokensPerSecond)

	return nil
}
