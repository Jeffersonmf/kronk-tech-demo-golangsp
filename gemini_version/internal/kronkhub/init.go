// Package kronkhub provides shared helpers for loading local GGUF models
// via the Kronk SDK.
package kronkhub

import (
	"fmt"

	"github.com/ardanlabs/kronk/sdk/kronk"
	"github.com/ardanlabs/kronk/sdk/kronk/model"
)

// Paths to the GGUF models already downloaded via `kronk model pull --local`.
const (
	ChatModelGPTOSS = "/home/jmarchetti/.kronk/models/unsloth/gpt-oss-20b-GGUF/gpt-oss-20b-Q8_0.gguf"
	ChatModelQwen3  = "/home/jmarchetti/.kronk/models/Qwen/Qwen3-8B-GGUF/Qwen3-8B-Q8_0.gguf"
	EmbedModel      = "/home/jmarchetti/.kronk/models/ggml-org/embeddinggemma-300m-qat-q8_0-GGUF/embeddinggemma-300m-qat-Q8_0.gguf"
	RerankModel     = "/home/jmarchetti/.kronk/models/gpustack/bge-reranker-v2-m3-GGUF/bge-reranker-v2-m3-Q8_0.gguf"
)

// NewModel initializes the Kronk runtime and loads the model at path.
func NewModel(path string, opts ...model.Option) (*kronk.Kronk, error) {
	if err := kronk.Init(); err != nil {
		return nil, fmt.Errorf("kronk init: %w", err)
	}

	base := []model.Option{model.WithModelFiles([]string{path})}

	krn, err := kronk.New(append(base, opts...)...)
	if err != nil {
		return nil, fmt.Errorf("kronk new: %w", err)
	}

	return krn, nil
}
