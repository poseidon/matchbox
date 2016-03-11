package api

import (
	"fmt"

	"github.com/coreos/coreos-baremetal/bootcfg/storage"
	"github.com/coreos/coreos-baremetal/bootcfg/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfile = &storagepb.Profile{
		Id: "g1h2i3j4",
		Boot: &storagepb.NetBoot{
			Kernel: "/image/kernel",
			Initrd: []string{"/image/initrd_a", "/image/initrd_b"},
			Cmdline: map[string]string{
				"a": "b",
				"c": "",
			},
		},
		CloudId:    "cloud-config.yml",
		IgnitionId: "ignition.json",
	}

	testProfileIgnitionYAML = &storagepb.Profile{
		Id:         "g1h2i3j4",
		IgnitionId: "ignition.yaml",
	}

	testGroup = &storagepb.Group{
		Id:           "test-group",
		Name:         "test group",
		Profile:      "g1h2i3j4",
		Requirements: map[string]string{"uuid": "a1b2c3d4"},
		Metadata:     []byte(`{"k8s_version":"v1.1.2","service_name":"etcd2","pod_network": "10.2.0.0/16"}`),
	}

	testGroupWithMAC = &storagepb.Group{
		Id:           "test-group",
		Name:         "test group",
		Profile:      "g1h2i3j4",
		Requirements: map[string]string{"mac": validMACStr},
	}
)

type fixedStore struct {
	Groups          map[string]*storagepb.Group
	Profiles        map[string]*storagepb.Profile
	IgnitionConfigs map[string]string
	CloudConfigs    map[string]string
}

func (s *fixedStore) GroupGet(id string) (*storagepb.Group, error) {
	if group, present := s.Groups[id]; present {
		return group, nil
	}
	return nil, storage.ErrGroupNotFound
}

func (s *fixedStore) GroupList() ([]*storagepb.Group, error) {
	groups := make([]*storagepb.Group, len(s.Groups))
	i := 0
	for _, g := range s.Groups {
		groups[i] = g
		i++
	}
	return groups, nil
}

func (s *fixedStore) ProfileGet(id string) (*storagepb.Profile, error) {
	if profile, present := s.Profiles[id]; present {
		return profile, nil
	}
	return nil, storage.ErrProfileNotFound
}

func (s *fixedStore) ProfileList() ([]*storagepb.Profile, error) {
	profiles := make([]*storagepb.Profile, len(s.Profiles))
	i := 0
	for _, p := range s.Profiles {
		profiles[i] = p
		i++
	}
	return profiles, nil
}

func (s *fixedStore) IgnitionGet(id string) (string, error) {
	if config, present := s.IgnitionConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Ignition Config %s", id)
}

func (s *fixedStore) CloudGet(id string) (string, error) {
	if config, present := s.CloudConfigs[id]; present {
		return config, nil
	}
	return "", fmt.Errorf("no Cloud Config %s", id)
}

type emptyStore struct{}

func (s *emptyStore) GroupGet(id string) (*storagepb.Group, error) {
	return nil, storage.ErrGroupNotFound
}

func (s *emptyStore) GroupList() (groups []*storagepb.Group, err error) {
	return groups, nil
}

func (s *emptyStore) ProfileGet(id string) (*storagepb.Profile, error) {
	return nil, storage.ErrProfileNotFound
}

func (s *emptyStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	return profiles, nil
}

func (s *emptyStore) IgnitionGet(id string) (string, error) {
	return "", fmt.Errorf("no Ignition Config %s", id)
}

func (s *emptyStore) CloudGet(id string) (string, error) {
	return "", fmt.Errorf("no Cloud Config %s", id)
}
