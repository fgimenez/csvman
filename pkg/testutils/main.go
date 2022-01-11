package testutils

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

func CreateFiles(appFs afero.Fs, files map[string]string) error {
	for path, content := range files {
		f, err := appFs.Create(path)

		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(content)

		if err != nil {
			return err
		}
	}
	return nil
}

func CheckError(err error, expected bool, msg string) error {
	if !expected && err != nil {
		return fmt.Errorf("Unexpected error %v", err)
	}
	if expected {
		if err == nil {
			return fmt.Errorf("Expected error didn't happen")
		}
		if msg == "" {
			return fmt.Errorf("Expected error message not defined")
		}
		if !strings.HasPrefix(err.Error(), msg) {
			return fmt.Errorf("Wrong error, expected %q, got %q", msg, err.Error())
		}
	}
	return nil
}
