package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	fake "github.com/coreos/matchbox/matchbox/storage/testfakes"
)

func TestGroupCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	// assert that:
	// - Group creation was successful
	// - Group can be retrieved by id
	// - Group can be deleted by id
	err = store.GroupPut(fake.Group)
	assert.Nil(t, err)

	group, err := store.GroupGet(fake.Group.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)

	err = store.GroupDelete(fake.Group.Id)
	assert.Nil(t, err)
	_, err = store.GroupGet(fake.Group.Id)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestGroupGet(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Groups: map[string]*storagepb.Group{
			fake.Group.Id:           fake.Group,
			fake.GroupNoMetadata.Id: fake.GroupNoMetadata,
		},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	// assert that:
	// - Groups written to the store can be retrieved
	group, err := store.GroupGet(fake.Group.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)
	group, err = store.GroupGet(fake.GroupNoMetadata.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.GroupNoMetadata, group)
}

func TestGroupGet_NoGroup(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	_, err = store.GroupGet("no-such-group")
	if assert.Error(t, err) {
		assert.IsType(t, &os.PathError{}, err)
	}
}

func TestGroupList(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Groups: map[string]*storagepb.Group{
			fake.Group.Id:           fake.Group,
			fake.GroupNoMetadata.Id: fake.GroupNoMetadata,
		},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	groups, err := store.GroupList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(groups)) {
		assert.Contains(t, groups, fake.Group)
		assert.Contains(t, groups, fake.GroupNoMetadata)
		assert.NotContains(t, groups, &storagepb.Group{})
	}
}

func TestProfileCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	// assert that:
	// - Profile creation was successful
	// - Profile can be retrieved by id
	// - Profile can be deleted by id
	err = store.ProfilePut(fake.Profile)
	assert.Nil(t, err)

	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Nil(t, err)
	assert.Equal(t, fake.Profile, profile)

	err = store.ProfileDelete(fake.Profile.Id)
	assert.Nil(t, err)
	_, err = store.ProfileGet(fake.Profile.Id)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestProfileGet(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Equal(t, fake.Profile, profile)
	assert.Nil(t, err)
	_, err = store.ProfileGet("no-such-profile")
	if assert.Error(t, err) {
		assert.IsType(t, &os.PathError{}, err)
	}
}

func TestProfileList(t *testing.T) {
	dir, err := setup(&fake.FixedStore{
		Profiles: map[string]*storagepb.Profile{fake.Profile.Id: fake.Profile},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	profiles, err := store.ProfileList()
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(profiles)) {
		assert.Equal(t, fake.Profile, profiles[0])
	}
}

func TestIgnitionCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	// assert that:
	// - Ignition template creation was successful
	// - Ignition template can be retrieved by name
	// - Ignition template can be deleted by name
	err = store.IgnitionPut(fake.IgnitionYAMLName, []byte(fake.IgnitionYAML))
	assert.Nil(t, err)

	template, err := store.IgnitionGet(fake.IgnitionYAMLName)
	assert.Nil(t, err)
	assert.Equal(t, fake.IgnitionYAML, template)

	err = store.IgnitionDelete(fake.IgnitionYAMLName)
	assert.Nil(t, err)
	_, err = store.IgnitionGet(fake.IgnitionYAMLName)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestIgnitionGet(t *testing.T) {
	contents := `{"ignitionVersion":1,"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`
	dir, err := setup(&fake.FixedStore{
		IgnitionConfigs: map[string]string{"myignition.json": contents},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	ign, err := store.IgnitionGet("myignition.json")
	assert.Equal(t, contents, ign)
	assert.Nil(t, err)
}

func TestGenericCRUD(t *testing.T) {
	dir, err := setup(&fake.FixedStore{})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	// assert that:
	// - Generic template creation was successful
	// - Generic template can be retrieved by name
	// - Generic template can be deleted by name
	err = store.GenericPut(fake.GenericName, []byte(fake.Generic))
	assert.Nil(t, err)

	template, err := store.GenericGet(fake.GenericName)
	assert.Nil(t, err)
	assert.Equal(t, fake.Generic, template)

	err = store.GenericDelete(fake.GenericName)
	assert.Nil(t, err)
	_, err = store.GenericGet(fake.GenericName)
	if assert.Error(t, err) {
		assert.IsType(t, err, &os.PathError{})
	}
}

func TestGenericGet(t *testing.T) {
	contents := `{"ignitionVersion":1,"storage":{},"systemd":{"units":[{"name":"etcd2.service","enable":true}]},"networkd":{},"passwd":{}}`
	dir, err := setup(&fake.FixedStore{
		GenericConfigs: map[string]string{"generic": contents},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	ign, err := store.GenericGet("generic")
	assert.Equal(t, contents, ign)
	assert.Nil(t, err)
}

func TestCloudGet(t *testing.T) {
	contents := "#cloud-config"
	dir, err := setup(&fake.FixedStore{
		CloudConfigs: map[string]string{"cloudcfg.yaml": contents},
	})
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	store := NewFileStore(&Config{Root: dir})
	cfg, err := store.CloudGet("cloudcfg.yaml")
	assert.Nil(t, err)
	assert.Equal(t, contents, cfg)
}

// setup creates a temp fileStore directory to mirror a given fixedStore
// for testing. Returns the directory tree root. The caller must remove the
// temp directory when finished.
func setup(fixedStore *fake.FixedStore) (root string, err error) {
	root, err = ioutil.TempDir("", "data")
	if err != nil {
		return "", err
	}
	// directories
	profileDir := filepath.Join(root, "profiles")
	groupDir := filepath.Join(root, "groups")
	ignitionDir := filepath.Join(root, "ignition")
	genericDir := filepath.Join(root, "generic")
	cloudDir := filepath.Join(root, "cloud")
	if err := mkdirs(profileDir, groupDir, ignitionDir, genericDir, cloudDir); err != nil {
		return root, err
	}
	// files
	for _, profile := range fixedStore.Profiles {
		profileFile := filepath.Join(profileDir, profile.Id+".json")
		data, err := json.MarshalIndent(profile, "", "\t")
		if err != nil {
			return root, err
		}
		err = ioutil.WriteFile(profileFile, []byte(data), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for _, group := range fixedStore.Groups {
		groupFile := filepath.Join(groupDir, group.Id+".json")
		richGroup, err := group.ToRichGroup()
		if err != nil {
			return root, err
		}
		data, err := json.MarshalIndent(richGroup, "", "\t")
		if err != nil {
			return root, err
		}
		err = ioutil.WriteFile(groupFile, []byte(data), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for name, content := range fixedStore.IgnitionConfigs {
		ignitionFile := filepath.Join(ignitionDir, name)
		err = ioutil.WriteFile(ignitionFile, []byte(content), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for name, content := range fixedStore.GenericConfigs {
		genericFile := filepath.Join(genericDir, name)
		err = ioutil.WriteFile(genericFile, []byte(content), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	for name, content := range fixedStore.CloudConfigs {
		cloudConfigFile := filepath.Join(cloudDir, name)
		err = ioutil.WriteFile(cloudConfigFile, []byte(content), defaultFileMode)
		if err != nil {
			return root, err
		}
	}
	return root, nil
}

// mkdirs creates new directories with the given names and default permission
// bits.
func mkdirs(names ...string) error {
	for _, dir := range names {
		if err := os.Mkdir(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}
	return nil
}
