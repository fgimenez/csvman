package main

import (
	"log"
	"os"

	"github.com/spf13/afero"

	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/runner"
)

func main() {
	fs := afero.NewOsFs()

	cfg, err := config.InitializeConfig(fs, os.Args[1:])
	if err != nil {
		log.Fatalf("Initializing config: %q", err)
	}

	err = runner.Transform(cfg)
	if err != nil {
		log.Fatalf("Executing transform: %q", err)
	}
}
