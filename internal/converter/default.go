package converter

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type DefaultConverter struct{}

func (dc *DefaultConverter) Supports(filePath string) bool {
	return true // Always supports as a fallback
}

func (dc *DefaultConverter) Convert(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	// --- Fast Initial Check (Header-Based) ---
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("error reading file header %s: %w", filePath, err)
	}

	contentType := http.DetectContentType(buffer)
	if strings.HasPrefix(contentType, "text/") || contentType == "application/json" || contentType == "application/xml" {
		// Likely text. Seek back and read the entire file.
		_, err = file.Seek(0, 0)
		if err != nil {
			return "", fmt.Errorf("error seeking to beginning of file %s: %w", filePath, err)
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("error reading %s: %w", filePath, err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	// --- Robust Binary Check (Non-Printable Characters) ---
	// If the header check didn't identify it as text, do a more thorough check.
	_, err = file.Seek(0, 0) // Reset file pointer.
	if err != nil {
		return "", fmt.Errorf("error seeking to beginning of file %s: %w", filePath, err)
	}
	content, err := os.ReadFile(filePath) // Read the whole file this time.
	if err != nil {
		return "", fmt.Errorf("error reading %s: %w", filePath, err)
	}

	if isBinary(content) {
		return fmt.Sprintf("Skipping binary file: %s", filePath), nil
	}

	return strings.TrimSpace(string(content)), nil
}
func isBinary(data []byte) bool {
	// Define a threshold for the percentage of non-printable characters.
	nonPrintableThreshold := 0.1 // 10%

	nonPrintableCount := 0
	for _, r := range string(data) { // Iterate over runes (Unicode code points).
		if r == unicode.ReplacementChar || (!unicode.IsPrint(r) && !unicode.IsSpace(r)) {
			nonPrintableCount++
		}
	}

	// Calculate the percentage of non-printable characters.
	nonPrintablePercentage := float64(nonPrintableCount) / float64(len(data))

	return nonPrintablePercentage > nonPrintableThreshold
}

func (dc *DefaultConverter) GetLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx", ".h", ".hpp", ".hxx":
		return "cpp"
	case ".c":
		return "c"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".sh":
		return "bash"
	case ".md":
		return "markdown"
	case ".lua":
		return "lua"
	case ".yml", ".yaml":
		return "yaml"
	case ".go":
		return "go"
	default:
		return ""
	}
}
