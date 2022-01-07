package rules_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/rules"
	"github.com/fgimenez/csvman/pkg/testutils"
)

func Test_Parse(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description    string
		files          map[string]string
		inputCfg       *config.Config
		expectedErr    bool
		expectedErrMsg string
		expectedRules  *rules.Rules
	}{
		{
			description: "yaml error",
			files: map[string]string{
				"/a/b/rules.yaml": "not a valid yaml",
			},
			inputCfg: &config.Config{
				RulesPath: "/a/b/rules.yaml",
				Fs:        appFs,
			},
			expectedErr:    true,
			expectedErrMsg: "yaml: unmarshal errors",
		},
		{
			description: "happy path",
			files: map[string]string{
				"/a/b/rules.yaml": `---
swapColumn:
- originIndex: 1
  targetIndex: 2
- originIndex: 3
  targetIndex: 4
renameColumn:
- index: 1
  name: newName1
- index: 2
  name: newName2
dropRow:
- regexp: pattern1
- regexp: pattern2
`,
			},
			inputCfg: &config.Config{
				RulesPath: "/a/b/rules.yaml",
				Fs:        appFs,
			},
			expectedRules: &rules.Rules{
				SwapColumn: []rules.SwapColumnRule{
					{
						OriginIndex: 1,
						TargetIndex: 2,
					},
					{
						OriginIndex: 3,
						TargetIndex: 4,
					},
				},
				RenameColumn: []rules.RenameColumnRule{
					{
						Index: 1,
						Name:  "newName1",
					},
					{
						Index: 2,
						Name:  "newName2",
					},
				},
				DropRow: []rules.DropRowRule{
					{
						Regexp: "pattern1",
					},
					{
						Regexp: "pattern2",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			err := testutils.CreateFiles(appFs, tc.files)
			if err != nil {
				t.Fatalf("could not create files %v", err)
			}

			actualRules, err := rules.Parse(tc.inputCfg)

			if !tc.expectedErr && err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if tc.expectedErr {
				if err == nil {
					t.Fatalf("Expected error didn't happen")
				}
				if tc.expectedErrMsg == "" {
					t.Fatalf("Expected error message not defined")
				}
				if !strings.HasPrefix(err.Error(), tc.expectedErrMsg) {
					t.Fatalf("Wrong error, expected %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}

			if !reflect.DeepEqual(actualRules, tc.expectedRules) {
				t.Fatalf("Actual rules %#v don't match expected rules %#v", actualRules, tc.expectedRules)
			}
		})
	}
}

func Test_Apply(t *testing.T) {

	tcs := []struct {
		description string
	}{
		{
			description: "",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {

		})
	}
}
