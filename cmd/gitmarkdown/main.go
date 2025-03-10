package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitmarkdown/internal/converter"
	"gitmarkdown/internal/exporter"
	"gitmarkdown/internal/tree"
	"gitmarkdown/pkg/utils"

	"github.com/gobwas/glob"
)

func main() {
	inputPath := flag.String("input", "./", "Path to file or directory")
	outputPath := flag.String("output", "", "Optional output file to save results")
	copyToClipboard := flag.Bool("copy", false, "Copy output to clipboard")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	ignore := flag.String("ignore", "", "Ignore specific files (comma-separated, supports wildcards, e.g., '*.css,assets/*.html')")

	flag.Parse()

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		log.Fatalf("Input path '%s' not found.", *inputPath)
	}

	convertersList := []converter.Converter{
		&converter.DefaultConverter{}, // Include DefaultConverter
	}
	exportersList := []exporter.Exporter{
		&exporter.MarkdownExporter{},
	}

	if len(convertersList) == 0 || len(exportersList) == 0 {
		log.Fatal("No converters or exporters found. Exiting.")
	}

	var defaultConverter converter.Converter = &converter.DefaultConverter{}

	exporter := exportersList[0]

	var ignorePatterns []string
	if *ignore != "" {
		ignorePatterns = strings.Split(*ignore, ",")
		for i := range ignorePatterns {
			ignorePatterns[i] = strings.TrimSpace(ignorePatterns[i])
		}
	}

	var output strings.Builder

	fileInfo, err := os.Stat(*inputPath)
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	if !fileInfo.IsDir() {
		conv := converter.GetConverter(*inputPath, convertersList, defaultConverter)
		content, err := conv.Convert(*inputPath)
		if err != nil {
			log.Fatalf("Error during conversion: %v", err)
		}
		language := conv.GetLanguage(*inputPath)
		output.WriteString(exporter.Format(*inputPath, content, language))
	} else {
		treeOutput, err := processDirectory(*inputPath, ignorePatterns)
		if err != nil {
			log.Fatalf("Error processing directory: %v", err) // Handle the error
		}

		output.WriteString(fmt.Sprintf("## Tree for %s\n```\n%s\n```\n\n", filepath.Base(*inputPath), treeOutput))

		err = filepath.WalkDir(*inputPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() && d.Name() == ".git" {
				return filepath.SkipDir
			}

			if !d.IsDir() {
				relPath, _ := filepath.Rel(*inputPath, path)

				// Compile ignore patterns to globs
				var globs []glob.Glob
				for _, pattern := range ignorePatterns {
					g, err := glob.Compile(pattern)
					if err != nil {
						log.Printf("Invalid glob pattern '%s': %v", pattern, err)
						continue // Skip invalid patterns
					}
					globs = append(globs, g)
				}

				// Check if the file should be ignored
				shouldIgnore := false
				for _, g := range globs {
					if g.Match(relPath) {
						shouldIgnore = true
						break
					}
				}

				globalIgnore := utils.LoadGlobalIgnorePatterns()
				if utils.ShouldIgnore(relPath, globalIgnore, *inputPath) {
					shouldIgnore = true
				}

				if shouldIgnore {
					return nil
				}

				fileInfo, err := d.Info()
				if err != nil {
					return err
				}

				if fileInfo.Size() == 0 {
					if *verbose {
						log.Printf("Skipping empty file: %s", path)
					}
					return nil
				}

				conv := converter.GetConverter(path, convertersList, defaultConverter)
				content, err := conv.Convert(path)
				if err != nil {
					return fmt.Errorf("error converting file %s: %w", path, err)
				}
				language := conv.GetLanguage(path)
				output.WriteString(exporter.Format(relPath, content, language))
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error during directory traversal: %v", err)
		}
	}

	if *outputPath != "" {
		err := os.WriteFile(*outputPath, []byte(output.String()), 0644)
		if err != nil {
			log.Fatalf("Error writing to output file: %v", err)
		}
		log.Printf("Output written to %s", *outputPath)
	}

	if *copyToClipboard {
		err := utils.CopyContent(output.String())
		if err != nil {
			log.Printf("Error copying to clipboard: %v", err)
		} else {
			log.Println("Output copied to clipboard")
		}
	}

	// If no output path or copy flag, print to stdout.
	if *outputPath == "" && !*copyToClipboard {
		fmt.Println(output.String())
	}
}

func processDirectory(directory string, ignorePatterns []string) (string, error) {
	treeData, err := tree.BuildTree(directory, ignorePatterns)
	if err != nil {
		return "", err // Propagate the error
	}
	return tree.FormatTree(treeData, ""), nil
}
