package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/afero"
)

type Config struct {
	Fs afero.Fs

	RulesPath  string
	InputPath  string
	OutputPath string
}

func InitializeConfig(fs afero.Fs, args []string) (*Config, error) {
	flags := flag.FlagSet{}
	cfg := &Config{
		Fs: fs,
	}

	flags.StringVar(&cfg.RulesPath, "rules", "rules.yaml", "Path to rules file")
	flags.StringVar(&cfg.InputPath, "input", "input.csv", "Path to input file")
	flags.StringVar(&cfg.OutputPath, "output", "output.csv", "Path to output file")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	err = validate(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if _, err := cfg.Fs.Stat(cfg.RulesPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Rules file path %q does't exist", cfg.RulesPath)
	}
	if _, err := cfg.Fs.Stat(cfg.InputPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Input file path %q does't exist", cfg.InputPath)
	}

	return nil
}
