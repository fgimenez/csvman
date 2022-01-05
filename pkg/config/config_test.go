package config_test

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
)

func Test_InitializeFlags(t *testing.T) {
	tcs := []struct {
		description string
		args        []string
		files       map[string]string
		expectedErr bool
		expectedCfg *config.Config
	}{
		{
			description: "invalid rules path gives error",
			args: []string{
				"-rules",
				"not a real path",
			},
			expectedErr: true,
			expectedCfg: nil,
		},
	}

	var appFs = afero.NewMemMapFs()

	for _, tc := range tcs {
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
			if tc.expectedErr && err == nil {
				t.Fatalf("Expected error didn't happen")
			}

			if !reflect.DeepEqual(actualCfg, tc.expectedCfg) {
				t.Fatalf("Actual flags %#v don't match expected flags %#v", actualCfg, tc.expectedCfg)
			}
		})
	}
}
