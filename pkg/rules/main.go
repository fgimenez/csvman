package rules

import (
	"encoding/csv"

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
	Regexp string `yaml:"regexp"`
}

type Rules struct {
	SwapColumn   []SwapColumnRule   `yaml:"swapColumn,omitempty"`
	RenameColumn []RenameColumnRule `yaml:"renameColumn,omitempty"`
	DropRow      []DropRowRule      `yaml:"dropRow,omitempty"`
}

func Parse(cfg *config.Config) (*Rules, error) {
	yamlFile, err := afero.ReadFile(cfg.Fs, cfg.RulesPath)
	if err != nil {
		return nil, err
	}
	r := &Rules{}
	err = yaml.Unmarshal(yamlFile, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rules) Apply(input csv.Reader, output csv.Writer) error {
	return nil
}
