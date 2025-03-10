# gitmarkdown: Convert Folders to LLM - ingestible markdown file

`gitmarkdown` is a command-line tool written in Go that converts the contents of a Git repository (or any directory) into a Markdown document.  It intelligently handles both code and directory structures, making it ideal for generating documentation, project overviews, or quick snapshots of your codebase.

## Features

*   **Code Highlighting:**  Automatically detects the programming language of each file (based on extension) and generates Markdown code blocks with appropriate syntax highlighting.  Supports a wide range of languages (Python, JavaScript, TypeScript, HTML, CSS, Java, C++, C, C#, Ruby, PHP, JSON, XML, Bash, Markdown, Lua, YAML, and Go).
*   **Directory Tree Generation:** Creates a nicely formatted text representation of the directory structure, making it easy to visualize the project layout.
*   **Binary File Handling:**  Robustly detects and skips binary files, preventing garbled output and keeping the focus on human-readable content.  This is done using a combination of header-based detection (for speed) and a character-based analysis (for accuracy).
*   **Ignore Files:**  Respects `.gitignore` and `.globalignore` files, allowing you to exclude unwanted files and directories from the output.  Also supports custom ignore patterns via the command line.
*   **Clipboard Support:**  Optionally copies the generated Markdown directly to your clipboard (works on Windows, macOS, and Linux with `xclip`, `xsel`, or `wl-copy`).
*   **Fast and Efficient:**  Written in Go for performance and minimal resource usage.  Handles large repositories efficiently.
*   **Easy to Use:**  Simple command-line interface with intuitive options.

## Installation

### Prerequisites

*   **Go:** You need Go (version 1.18 or later) installed and configured on your system.  See [the official Go installation instructions](https://go.dev/doc/install).
* **(Optional, for Linux clipboard)**: If you want to use the `-copy` option on Linux, you'll need either `xclip`, `xsel`, or `wl-copy` installed.
   ```bash
   # For xclip (X11)
   sudo apt-get update  # Debian/Ubuntu
   sudo apt-get install xclip

   sudo yum install xclip # Fedora/CentOS/RHEL

   # For xsel (X11)
    sudo apt-get update  # Debian/Ubuntu
    sudo apt-get install xsel

    sudo yum install xsel # Fedora/CentOS/RHEL

   # For wl-copy (Wayland)
   sudo apt-get update  #Debian/Ubuntu
   sudo apt-get install wl-clipboard

   sudo yum install wl-clipboard #Fedora/CentOS/RHEL
   ```

### From Binaries

Download the pre-built binaries for your operating system from the [Releases](https://github.com/pranjalya/gitmarkdown/releases) page.  Extract the archive and place the `gitmarkdown` executable in a directory included in your system's `PATH`.

### Using Homebrew (macOS)

```bash
brew install pranjalya/tap/gitmarkdown
```

### Using Go

```bash
go install github.com/pranjalya/gitmarkdown@latest
```

### Build from Source (Recommended)

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/Pranjalya/gitmarkdown.git
    cd gitmarkdown
    ```
2.  **Initialize the Go module:**

    ```bash
    go mod init gitmarkdown
    ```

3.  **Install dependencies:**
    ```bash
    go get github.com/gobwas/glob
    ```

4.  **Build the executable:**

    ```bash
    mkdir -p build
    go build -o build/gitmarkdown ./cmd/gitmarkdown
    ```

    This will create an executable named `gitmarkdown` (or `gitmarkdown.exe` on Windows) in the `build` directory.

5.  **(Optional) Move the executable:**  You can move the `gitmarkdown` executable to a directory in your `PATH` for easy access from anywhere.  For example:

    ```bash
    sudo mv buiuld/gitmarkdown /usr/local/bin/
    ```

## Usage

```bash
gitmarkdown [options]
```

**Options:**

*   **`-input <path>`:**  The path to the file or directory you want to convert.  Defaults to the current directory (`.`).
*   **`-output <file>`:**  The path to the output Markdown file.  If not specified, the output is printed to the console (stdout).
*   **`-copy`:**  Copy the generated Markdown to the clipboard.
*   **`-verbose`:**  Enable verbose logging (useful for debugging).
*   **`-ignore <pattern1,pattern2,...>`:**  A comma-separated list of file/directory patterns to ignore (supports wildcards, e.g., `*.log,temp/*`).  These patterns are *in addition* to those specified in `.gitignore` and `.globalignore`.

**Examples:**

*   Convert the current directory and print the output to the console:

    ```bash
    gitmarkdown
    ```

*   Convert a specific directory and save the output to `output.md`:

    ```bash
    gitmarkdown -input my_project -output output.md
    ```

*   Convert a single file and copy the output to the clipboard:

    ```bash
    gitmarkdown -input my_file.py -copy
    ```

*   Convert a directory, ignoring `.log` files and the `temp` directory, and save to `docs.md`:

    ```bash
    gitmarkdown -input my_project -ignore "*.log,temp/*" -output docs.md
    ```

*  Convert a directory with verbose:
    ```bash
    gitmarkdown -input my_project -verbose
    ```
## How it Works

1.  **Directory Traversal:**  `gitmarkdown` recursively walks the input directory (or processes a single file if specified).
2.  **Ignore File Handling:**  It reads `.gitignore` and `.globalignore` files to determine which files and directories to exclude.  Command-line ignore patterns are also applied.
3.  **File Type Detection:** For each file:
    *   **Fast Header Check:** It first checks the file header (first 512 bytes) to quickly identify common text file types (using `net/http.DetectContentType`).
    *   **Robust Binary Check:** If the header check is inconclusive, it performs a more thorough analysis, counting non-printable Unicode characters.  If a threshold is exceeded, the file is considered binary and skipped.
4.  **Content Conversion:** For text files, the content is read and any leading/trailing whitespace is trimmed.
5.  **Language Detection:** The programming language is determined based on the file extension.
6.  **Markdown Formatting:**  The output is formatted as Markdown:
    *   A directory tree is generated for directory inputs.
    *   Each file's content is enclosed in a Markdown code block with appropriate syntax highlighting based on the detected language.
7.  **Output:** The generated Markdown is either written to the specified output file, copied to the clipboard, or printed to the console.

## Contributing

Contributions are welcome!  If you find a bug or have a feature request, please open an issue on GitHub.  If you'd like to contribute code, please fork the repository and submit a pull request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.