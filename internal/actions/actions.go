package actions

import (
	"fmt"
	"os"
)

const (
	fileCMDOutput = "GITHUB_OUTPUT"
)

func SetOutput(key string, value any) {
	err := fileCMD(fileCMDOutput, fmt.Sprintf("%s=%v", key, value))
	if err != nil {
		// We can't recover from this state
		panic(err)
	}
}

func fileCMD(cmd, msg string) error {
	filepath := os.Getenv(cmd)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("can not open environment file: %w", err)
	}

	defer f.Close()

	_, err = fmt.Fprintln(f, msg)
	if err != nil {
		return fmt.Errorf("can not write to environment file: %w", err)
	}

	return nil
}
