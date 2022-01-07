package runner

import (
	"github.com/fgimenez/csvman/pkg/config"
	"github.com/fgimenez/csvman/pkg/rules"
)

func Transform(cfg *config.Config) error {
	_, err := rules.Parse(cfg)
	if err != nil {
		return err
	}

	/*
		_, err = r.Apply(inputData)
		if err != nil {
			return err
		}
	*/

	return nil
}
