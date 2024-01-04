package tailo

import (
	"fmt"
	"os"
	"os/exec"
)

// Build runs the Tailwind CSS CLI binary to build the
// CSS file and generate compiled CSS it expects to find
// the options in the config file.
func Build(options ...Option) {
	// Applying passed options
	for _, option := range options {
		option()
	}

	fullPath, _, err := GetFullBinaryPath()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := Setup()
		if err != nil {
			panic(err)
		}
	}

	cmd := exec.Command(fullPath)
	cmd.Args = append(cmd.Args, "--config", configPath)
	cmd.Args = append(cmd.Args, "--input", inputPath)
	cmd.Args = append(cmd.Args, "--output", outputPath)

	cmd.Stdout = writer(0)
	cmd.Stderr = writer(0)

	fmt.Println("[tailo] Running:", cmd.String())

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
