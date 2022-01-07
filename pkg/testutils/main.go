package testutils

import "github.com/spf13/afero"

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
