package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/kronk/sdk/kronk/model"

	"github.com/jeff/gemini-kronk-demo/internal/kronkhub"
)

type target struct {
	name string
	path string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	targets := []target{
		{"gpt-oss-20b-Q8_0", kronkhub.ChatModelGPTOSS},
		{"Qwen3-8B-Q8_0", kronkhub.ChatModelQwen3},
		{"embeddinggemma-300m-qat-Q8_0", kronkhub.EmbedModel},
		{"bge-reranker-v2-m3-Q8_0", kronkhub.RerankModel},
	}

	for _, t := range targets {
		if err := checkModel(t); err != nil {
			return fmt.Errorf("%s: %w", t.name, err)
		}
	}

	return nil
}

func checkModel(t target) error {
	fmt.Printf("\n=== %s ===\n", t.name)

	start := time.Now()
	krn, err := kronkhub.NewModel(t.path)
	if err != nil {
		return err
	}
	loadDur := time.Since(start)

	defer func() {
		if err := krn.Unload(context.Background()); err != nil {
			fmt.Printf("unload error: %v\n", err)
		}
	}()

	info := krn.ModelInfo()
	fmt.Printf("load time   : %s\n", loadDur)
	fmt.Printf("IsGPTModel  : %v\n", info.IsGPTModel)
	fmt.Printf("IsEmbedModel: %v\n", info.IsEmbedModel)
	fmt.Printf("template    : %s\n", info.Template.FileName)

	if info.IsEmbedModel {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		resp, err := krn.Embeddings(ctx, model.D{"input": "test"})
		if err != nil {
			return fmt.Errorf("embeddings: %w", err)
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("embeddings: empty response data")
		}
		fmt.Printf("OK: embedding dimension = %d\n", len(resp.Data[0].Embedding))
	} else {
		fmt.Println("OK")
	}

	return nil
}
