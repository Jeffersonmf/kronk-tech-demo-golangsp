package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/kronk/sdk/kronk/model"

	"github.com/jeff/gemini-kronk-demo/internal/kronkhub"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("Starting Kronk Local Efficiency Hub...")
	fmt.Println("Mode: Local-Only (CPU Inference)")
	fmt.Println("Model:", kronkhub.ChatModelGPTOSS)

	krn, err := kronkhub.NewModel(kronkhub.ChatModelGPTOSS, model.WithContextWindow(8192))
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("\nUnloading model...")
		if err := krn.Unload(context.Background()); err != nil {
			fmt.Printf("failed to unload model: %v\n", err)
		}
	}()

	fmt.Print("- system info:\n\t")
	for k, v := range krn.SystemInfo() {
		fmt.Printf("%s:%v, ", k, v)
	}
	fmt.Println()

	fmt.Println("- contextWindow:", krn.ModelConfig().ContextWindow())
	fmt.Println("- embeddings   :", krn.ModelInfo().IsEmbedModel)
	fmt.Println("- isGPT        :", krn.ModelInfo().IsGPTModel)
	fmt.Println("- template     :", krn.ModelInfo().Template.FileName)

	fmt.Println("\nReady. Press Ctrl+C to stop.")

	<-ctx.Done()
	fmt.Println("\nStopping hub...")
	return nil
}
