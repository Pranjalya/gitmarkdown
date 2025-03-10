package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CopyContent(content string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows uses UTF-16LE encoding.
		utf16le, err := EncodeUTF16([]byte(content), true)
		if err != nil {
			return fmt.Errorf("error encoding to UTF-16LE: %w", err)
		}

		// Pipe the UTF-16LE encoded content to clip.
		cmd = exec.Command("cmd", "/c", "clip")
		cmd.Stdin = bytes.NewReader(utf16le)
		cmd.Stdout = os.Stdout // Redirect stdout for debugging.
		cmd.Stderr = os.Stderr // Redirect stderr for debugging.

	case "darwin":
		cmd = exec.Command("pbcopy")
		cmd.Stdin = bytes.NewBufferString(content)
	case "linux":
		// Check for Wayland or X11
		sessionType := os.Getenv("XDG_SESSION_TYPE")
		wayland := sessionType == "wayland" || os.Getenv("WAYLAND_DISPLAY") != ""
		//x11 := sessionType == "x11" || os.Getenv("DISPLAY") != ""

		if wayland {
			cmd = exec.Command("wl-copy")
			cmd.Stdin = bytes.NewBufferString(content)

		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard") // Try xclip first
			cmd.Stdin = bytes.NewBufferString(content)
		}
	default:
		return fmt.Errorf("clipboard not supported on this platform")
	}

	if err := cmd.Run(); err != nil {
		// If xclip failed, try xsel
		if runtime.GOOS == "linux" && cmd.Args[0] == "xclip" {
			cmd = exec.Command("xsel", "--clipboard", "--input")
			cmd.Stdin = bytes.NewBufferString(content)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error copying to clipboard using xsel: %w", err)
			}
		} else {

			return fmt.Errorf("error copying to clipboard: %w", err)
		}
	}
	return nil
}

// EncodeUTF16 encodes a byte slice to UTF-16 little-endian or big-endian.
func EncodeUTF16(b []byte, littleEndian bool) ([]byte, error) {
	if len(b)%2 != 0 {
		return nil, fmt.Errorf("input byte slice length must be even, got %d", len(b))
	}
	// Check for BOM (Byte Order Mark)
	hasBOM := false
	if len(b) >= 2 {
		if (b[0] == 0xFE && b[1] == 0xFF) || (b[0] == 0xFF && b[1] == 0xFE) {
			hasBOM = true
		}
	}

	var result []byte
	if littleEndian {
		// Add BOM if it doesn't exist
		if !hasBOM {
			result = append(result, 0xFF, 0xFE) // Little Endian BOM
		}

		for i := 0; i < len(b); i += 2 {
			// Swap bytes for little-endian
			if hasBOM { // if source has BOM - just add it
				result = append(result, b[i+1], b[i])
			} else {
				result = append(result, b[i], b[i+1])
			}

		}
	} else {
		// Add BOM if it doesn't exist
		if !hasBOM {
			result = append(result, 0xFE, 0xFF) // Big Endian BOM
		}

		for i := 0; i < len(b); i += 2 {
			// Keep bytes as is for big-endian
			result = append(result, b[i], b[i+1])
		}
	}
	return result, nil
}
