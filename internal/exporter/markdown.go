package exporter

import (
	"fmt"
)

type MarkdownExporter struct{}

func (me *MarkdownExporter) Format(relativePath, content, language string) string {
	if language != "" {
		return fmt.Sprintf("## File: %s\n```%s\n%s\n```\n", relativePath, language, content)
	}
	return fmt.Sprintf("## File: %s\n```\n%s\n```\n", relativePath, content)
}
