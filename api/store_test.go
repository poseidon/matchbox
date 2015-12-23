package api

import (
	"fmt"
)

type fixedStore struct {
	BootCfg  *BootConfig
	CloudCfg *CloudConfig
}

func (s *fixedStore) BootConfig(attrs MachineAttrs) (*BootConfig, error) {
	return s.BootCfg, nil
}

func (s fixedStore) CloudConfig(attrs MachineAttrs) (*CloudConfig, error) {
	return s.CloudCfg, nil
}

type emptyStore struct{}

func (s *emptyStore) BootConfig(attrs MachineAttrs) (*BootConfig, error) {
	return nil, fmt.Errorf("no boot config for machine")
}

func (s emptyStore) CloudConfig(attrs MachineAttrs) (*CloudConfig, error) {
	return nil, fmt.Errorf("no cloud config for machine")
}
