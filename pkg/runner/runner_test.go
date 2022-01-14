package runner_test

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/runner"
	"github.com/fgimenez/csvman/pkg/testutils"
)

const (
	rulesFile  = "/a/b/rules.yaml"
	inputFile  = "/a/b/input.csv"
	outputFile = "/a/b/output.csv"
)

var appFs = afero.NewMemMapFs()

func Test_Transform(t *testing.T) {
	tcs := []struct {
		description       string
		files             map[string]string
		inputCfg          *config.Config
		expectedErr       bool
		expectedErrMsg    string
		expectedOutputCsv string
	}{
		{
			description: "happy path",
			files: map[string]string{
				rulesFile: `---
swapColumn:
- originIndex: 0
  targetIndex: 2
renameColumn:
- index: 1
  name: "newField1"
dropRow:
- substring: "value11"
`,
				inputFile: `field0,field1,field2
value00,value01,value02
value10,value11,value12
value20,value21,value22`,
			},
			inputCfg: &config.Config{
				RulesPath:  rulesFile,
				InputPath:  inputFile,
				OutputPath: outputFile,
				Fs:         appFs,
			},
			expectedOutputCsv: `field2,newField1,field0
value02,value01,value00
value22,value21,value20
`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			err := testutils.CreateFiles(appFs, tc.files)
			if err != nil {
				t.Fatalf("could not create files %v", err)
			}

			err = runner.Transform(tc.inputCfg)
			if err := testutils.CheckError(err, tc.expectedErr, tc.expectedErrMsg); err != nil {
				t.Fatalf(err.Error())
			}

			afs := &afero.Afero{Fs: appFs}
			actualOutputCsvBytes, err := afs.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Unexpected error reading output file contents %v\n", err)
			}

			actualOutputCsv := string(actualOutputCsvBytes)

			if !reflect.DeepEqual(actualOutputCsv, tc.expectedOutputCsv) {
				t.Fatalf("Actual output CSV %#v doesn't match expected output CSV %#v", actualOutputCsv, tc.expectedOutputCsv)
			}
		})
	}

}
