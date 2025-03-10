package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitmarkdown/pkg/utils"
)

type TreeNode struct {
	Name     string
	IsDir    bool
	Children map[string]*TreeNode
}

func BuildTree(directory string, ignorePatterns []string) (*TreeNode, error) {
	node := &TreeNode{
		Name:     filepath.Base(directory),
		IsDir:    true,
		Children: make(map[string]*TreeNode),
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", directory, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(directory, entry.Name())
		relPath, _ := filepath.Rel(directory, fullPath) // No error check needed; it comes from ReadDir

		// Check for .git exclusion and ignore patterns
		if entry.IsDir() && entry.Name() == ".git" {
			continue
		}
		globalIgnore := utils.LoadGlobalIgnorePatterns()
		if utils.ShouldIgnore(relPath, globalIgnore, directory) {
			continue
		}

		if entry.IsDir() {
			childNode, err := BuildTree(fullPath, ignorePatterns)
			if err != nil {
				return nil, err
			}
			node.Children[entry.Name()] = childNode
		} else {
			node.Children[entry.Name()] = &TreeNode{Name: entry.Name(), IsDir: false}
		}
	}

	return node, nil
}

func FormatTree(node *TreeNode, prefix string) string {
	var lines []string
	if node.Children != nil {
		entries := make([]string, 0, len(node.Children))
		for name := range node.Children {
			entries = append(entries, name)
		}

		for i, name := range entries {
			child := node.Children[name]
			connector := "├── "
			if i == len(entries)-1 {
				connector = "└── "
			}
			line := prefix + connector + name
			if child.IsDir {
				line += "/"
			}
			lines = append(lines, line)

			if child.IsDir && child.Children != nil {
				extension := "│   "
				if i == len(entries)-1 {
					extension = "    "
				}
				deeper := FormatTree(child, prefix+extension)
				lines = append(lines, deeper)
			}
		}
	}
	return strings.Join(lines, "\n")
}
