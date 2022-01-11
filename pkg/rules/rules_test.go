package rules_test

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/rules"
	"github.com/fgimenez/csvman/pkg/testutils"
)

var appFs = afero.NewMemMapFs()

func Test_Parse(t *testing.T) {
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
- substring: pattern1
- substring: pattern2
`,
			},
			inputCfg: &config.Config{
				RulesPath: "/a/b/rules.yaml",
				Fs:        appFs,
			},
			expectedRules: &rules.Rules{
				Fs: appFs,
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
						Substring: "pattern1",
					},
					{
						Substring: "pattern2",
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
			if err := testutils.CheckError(err, tc.expectedErr, tc.expectedErrMsg); err != nil {
				t.Fatalf(err.Error())
			}

			if !reflect.DeepEqual(actualRules, tc.expectedRules) {
				t.Fatalf("Actual rules %#v don't match expected rules %#v", actualRules, tc.expectedRules)
			}
		})
	}
}

func Test_Apply(t *testing.T) {

	tcs := []struct {
		description       string
		inputRules        *rules.Rules
		inputCsv          string
		expectedErr       bool
		expectedErrMsg    string
		expectedOutputCsv string
	}{
		{
			description:    "swap column origin index out of bounds",
			expectedErr:    true,
			expectedErrMsg: "Swap column origin index out of bounds",
			inputCsv: `field0,field1
value0,value1
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				SwapColumn: []rules.SwapColumnRule{
					{
						OriginIndex: 10,
					},
				},
			},
		},
		{
			description:    "swap column target index out of bounds",
			expectedErr:    true,
			expectedErrMsg: "Swap column target index out of bounds",
			inputCsv: `field0,field1
value0,value1
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				SwapColumn: []rules.SwapColumnRule{
					{
						TargetIndex: 10,
					},
				},
			},
		},
		{
			description: "swap column happy path",
			inputCsv: `field0,field1,field2
value00,value01,value02
value10,value11,value12
value20,value21,value22
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				SwapColumn: []rules.SwapColumnRule{
					{
						OriginIndex: 0,
						TargetIndex: 2,
					},
				},
			},
			expectedOutputCsv: `field2,field1,field0
value02,value01,value00
value12,value11,value10
value22,value21,value20
`,
		},
		{
			description:    "rename column index out of bounds",
			expectedErr:    true,
			expectedErrMsg: "Rename column index out of bounds",
			inputCsv: `field0,field1
value0,value1
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				RenameColumn: []rules.RenameColumnRule{
					{
						Index: 10,
					},
				},
			},
		},
		{
			description:    "rename column empty name",
			expectedErr:    true,
			expectedErrMsg: "Rename column empty name",
			inputCsv: `field0,field1
value0,value1
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				RenameColumn: []rules.RenameColumnRule{
					{
						Name: "",
					},
				},
			},
		},
		{
			description: "rename column happy path",
			inputCsv: `field0,field1,field2
value00,value01,value02
value10,value11,value12
value20,value21,value22
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				RenameColumn: []rules.RenameColumnRule{
					{
						Index: 0,
						Name:  "newName",
					},
				},
			},
			expectedOutputCsv: `newName,field1,field2
value00,value01,value02
value10,value11,value12
value20,value21,value22
`,
		},
		{
			description: "drop row happy path",
			inputCsv: `field0,field1,field2
value00,value01,value02
value10,value_somestring_11,value12
value20,value21,value22
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				DropRow: []rules.DropRowRule{
					{
						Substring: "somestring",
					},
				},
			},
			expectedOutputCsv: `field0,field1,field2
value00,value01,value02
value20,value21,value22
`,
		},
		{
			description: "complete example",
			inputCsv: `field0,field1,field2
value00,value01,value02
value10,value_somestring_11,value12
value20,value21,value22
value30,value_someotherstring_31,value32
`,
			inputRules: &rules.Rules{
				Fs: appFs,
				SwapColumn: []rules.SwapColumnRule{
					{
						OriginIndex: 0,
						TargetIndex: 2,
					},
				},
				RenameColumn: []rules.RenameColumnRule{
					{
						Index: 0,
						Name:  "newName0",
					},
					{
						Index: 1,
						Name:  "newName1",
					},
				},
				DropRow: []rules.DropRowRule{
					{
						Substring: "somestring",
					},
					{
						Substring: "someotherstring",
					},
				},
			},
			expectedOutputCsv: `newName0,newName1,field0
value02,value01,value00
value22,value21,value20
`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			reader := csv.NewReader(strings.NewReader(tc.inputCsv))

			var output bytes.Buffer
			bufWriter := bufio.NewWriter(&output)
			writer := csv.NewWriter(bufWriter)

			err := tc.inputRules.Apply(*reader, *writer)

			if err := testutils.CheckError(err, tc.expectedErr, tc.expectedErrMsg); err != nil {
				t.Fatalf(err.Error())
			}

			actualOutputCsv := output.String()
			if actualOutputCsv != tc.expectedOutputCsv {
				t.Fatalf("Actual csv %#v don't match expected csv %#v", actualOutputCsv, tc.expectedOutputCsv)
			}
		})
	}
}
