// tailo is a wrapper for the Tailwind CSS CLI that
// facilitates the download and of the CLI and the
// config file.
package tailo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/labstack/gommon/log"
)

var (
	binaries = map[string]string{
		"darwin_amd64":  "tailwindcss-macos-x64",
		"darwin_arm64":  "tailwindcss-macos-arm64",
		"linux_amd64":   "tailwindcss-linux-x64",
		"linux_arm64":   "tailwindcss-linux-arm64",
		"linux_arm":     "tailwindcss-linux-armv7",
		"windows_amd64": "tailwindcss-windows-x64.exe",   // added windows endings
		"windows_arm64": "tailwindcss-windows-arm64.exe", // added windows endings
	}
)

// Setup downloads the Tailwind CSS CLI binary for the
// given operating system and architecture. It makes the
// binary executable and places it in the bin/ directory.
func Setup() error {

	fullPath, binary, err := GetFullBinaryPath()
	if err != nil {
		return fmt.Errorf("setup failed")
	}

	log.Printf("using filepath: %s", fullPath)

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Println("tailwind CSS CLI binary already exists.")

		return nil
	}

	url := fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/latest/download/%v", binary)
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	log.Printf("downloaded: %s", binary)

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("could not download Tailwind CSS CLI binary: %s", resp.Status)
	}

	// Create the path
	err = os.MkdirAll(binaryPath, 0755)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(fullPath, 0755)
	if err != nil {
		return err
	}

	log.Printf("tailwind can now be found at: %s", fullPath)

	return err
}

func GetFullBinaryPath() (string, string, error) {
	currentArch := fmt.Sprintf("%v_%v", runtime.GOOS, runtime.GOARCH)
	binary, ok := binaries[currentArch]
	if !ok {
		return "", "", fmt.Errorf("unsupported operating system and architecture: %s", currentArch)
	}
	fullPath := filepath.Join(binaryPath, binary)
	return fullPath, binary, nil
}
