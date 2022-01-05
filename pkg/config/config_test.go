package config_test

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
)

func Test_InitializeFlags(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description    string
		args           []string
		files          map[string]string
		expectedErr    bool
		expectedErrMsg string
		expectedCfg    *config.Config
	}{
		{
			description: "invalid rules path gives error",
			args: []string{
				"-input",
				"/a/b/input.csv",
				"-rules",
				"not a real path",
			},
			files: map[string]string{
				"/a/b/input.csv": "input content",
			},
			expectedErr:    true,
			expectedErrMsg: "Rules file path \"not a real path\" does't exist",
			expectedCfg:    nil,
		},
		{
			description: "invalid input path gives error",
			args: []string{
				"-input",
				"not a real path",
				"-rules",
				"/a/b/rules.csv",
			},
			files: map[string]string{
				"/a/b/rules.csv": "rules content",
			},
			expectedErr:    true,
			expectedErrMsg: "Input file path \"not a real path\" does't exist",
			expectedCfg:    nil,
		},
		{
			description: "valid input and rules paths give correct result",
			args: []string{
				"-input",
				"/a/b/input.csv",
				"-rules",
				"/a/b/rules.csv",
			},
			files: map[string]string{
				"/a/b/rules.csv": "rules content",
				"/a/b/input.csv": "input content",
			},
			expectedCfg: &config.Config{
				Fs:         appFs,
				RulesPath:  "/a/b/rules.csv",
				InputPath:  "/a/b/input.csv",
				OutputPath: "output.csv",
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			for path, content := range tc.files {
				f, err := appFs.Create(path)

				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}
				defer f.Close()

				_, err = f.WriteString(content)

				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}
			}

			actualCfg, err := config.InitializeConfig(appFs, tc.args)
			if !tc.expectedErr && err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if tc.expectedErr {
				if err == nil {
					t.Fatalf("Expected error didn't happen")
				}
				if tc.expectedErrMsg != err.Error() {
					t.Fatalf("Wrong error, expected %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}

			if !reflect.DeepEqual(actualCfg, tc.expectedCfg) {
				t.Fatalf("Actual flags %#v don't match expected flags %#v", actualCfg, tc.expectedCfg)
			}
		})
	}
}
