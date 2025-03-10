package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

var globalIgnoreFile = ".globalignore"

func LoadGlobalIgnorePatterns() []string {
	var patterns []string
	// Get the absolute path to .globalignore
	absGlobalIgnorePath, err := filepath.Abs(globalIgnoreFile)
	if err != nil {
		// Handle the error (e.g., log it or return an empty slice)
		return patterns // Return empty slice if there is an error getting abs path
	}
	file, err := os.Open(absGlobalIgnorePath)
	if err != nil {
		// Handle the error if the file doesn't exist or cannot be opened
		return patterns // Return empty if .globalignore not found.
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	// Check for any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		return patterns // Return what we have, even if incomplete.
	}

	return patterns
}

func ShouldIgnore(relativePath string, ignorePatterns []string, inputPath string) bool {

	// Always ignore .git directory
	if relativePath == ".git" || strings.HasPrefix(relativePath, ".git"+string(filepath.Separator)) {
		return true
	}

	// Load .gitignore patterns from the input directory
	gitIgnorePatterns := LoadGitIgnorePatterns(inputPath)

	// Combine .globalignore and .gitignore patterns
	combinedPatterns := append(ignorePatterns, gitIgnorePatterns...)

	// Compile ignore patterns into globs
	var globs []glob.Glob
	for _, pattern := range combinedPatterns {
		g, err := glob.Compile(pattern)
		if err != nil {
			// Optionally log or handle invalid patterns
			continue
		}
		globs = append(globs, g)
	}

	// Check if the path matches any of the glob patterns
	for _, g := range globs {
		if g.Match(relativePath) {
			return true
		}
	}

	return false
}

func LoadGitIgnorePatterns(dir string) []string {
	var patterns []string
	gitIgnorePath := filepath.Join(dir, ".gitignore")

	file, err := os.Open(gitIgnorePath)
	if err != nil {
		// If .gitignore doesn't exist, just return empty slice.
		return patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignore comments and empty lines
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns
}
