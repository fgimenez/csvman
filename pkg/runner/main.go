package runner

import (
	"encoding/csv"
	"os"

	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/rules"
)

func Transform(cfg *config.Config) error {
	r, err := rules.Parse(cfg)
	if err != nil {
		return err
	}

	inputFile, err := cfg.Fs.Open(cfg.InputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	input := csv.NewReader(inputFile)

	outputFile, err := cfg.Fs.OpenFile(cfg.OutputPath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	output := csv.NewWriter(outputFile)

	err = r.Apply(*input, *output)
	if err != nil {
		return err
	}

	return nil
}
