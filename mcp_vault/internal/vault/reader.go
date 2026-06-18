// Package vault reads Markdown documents from an Obsidian vault directory.
package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Document is a single Markdown file loaded from the vault.
type Document struct {
	// Path is relative to the vault root, e.g. "tasks/TASK-001.md".
	Path string
	// Content is the raw Markdown content of the file.
	Content string
}

// LoadAll walks root and returns every ".md" file found, sorted by path.
func LoadAll(root string) ([]Document, error) {
	var docs []Document

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("relpath %s: %w", path, err)
		}

		docs = append(docs, Document{Path: rel, Content: string(content)})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", root, err)
	}

	return docs, nil
}
