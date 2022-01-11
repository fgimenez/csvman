package rules

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"

	"github.com/fgimenez/csvman/pkg/config"
)

type SwapColumnRule struct {
	OriginIndex int `yaml:"originIndex"`
	TargetIndex int `yaml:"targetIndex"`
}

type RenameColumnRule struct {
	Index int    `yaml:"index"`
	Name  string `yaml:"name"`
}

type DropRowRule struct {
	Substring string `yaml:"substring"`
}

type Rules struct {
	Fs           afero.Fs
	SwapColumn   []SwapColumnRule   `yaml:"swapColumn,omitempty"`
	RenameColumn []RenameColumnRule `yaml:"renameColumn,omitempty"`
	DropRow      []DropRowRule      `yaml:"dropRow,omitempty"`
}

func Parse(cfg *config.Config) (*Rules, error) {
	yamlFile, err := afero.ReadFile(cfg.Fs, cfg.RulesPath)
	if err != nil {
		return nil, err
	}
	r := &Rules{Fs: cfg.Fs}
	err = yaml.Unmarshal(yamlFile, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rules) Apply(input csv.Reader, output csv.Writer) error {
	row := 0
	writeRecord := true
	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, rule := range r.SwapColumn {
			if len(record) < rule.OriginIndex {
				return fmt.Errorf("Swap column origin index out of bounds %d\n", rule.OriginIndex)
			}
			if len(record) < rule.TargetIndex {
				return fmt.Errorf("Swap column target index out of bounds %d\n", rule.TargetIndex)
			}

			record[rule.TargetIndex], record[rule.OriginIndex] = record[rule.OriginIndex], record[rule.TargetIndex]
		}

		for _, rule := range r.RenameColumn {
			if len(record) < rule.Index {
				return fmt.Errorf("Rename column index out of bounds %d\n", rule.Index)
			}

			if rule.Name == "" {
				return fmt.Errorf("Rename column empty name")
			}

			if row == 0 {
				record[rule.Index] = rule.Name
			}
		}

		for _, rule := range r.DropRow {
			if strings.Contains(strings.Join(record, ","), rule.Substring) {
				writeRecord = false
				continue
			}
		}

		if writeRecord {
			err = output.Write(record)
			if err != nil {
				return err
			}
		}
		writeRecord = true
		row++
	}
	output.Flush()
	return output.Error()
}
