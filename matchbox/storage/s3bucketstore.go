package storage

import (
	"encoding/json"
	"strings"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
	"github.com/sirupsen/logrus"
)

// S3BucketConfig initializes a fileStore.
type S3BucketConfig struct {
	BucketName string
	Region     string
	Logger     *logrus.Logger
}

// s3Store implements ths Store interface. Queries to the s3 bucket
// are restricted to the specified directory tree.
type s3BucketStore struct {
	client *S3Client
	logger *logrus.Logger
}

// NewS3BucketStore returns a new s3 bucket backed Store.
func NewS3BucketStore(config *S3BucketConfig) Store {
	client, _ := NewS3Client(config.Region, config.BucketName)

	return &s3BucketStore{
		client: client,
		logger: config.Logger,
	}
}

// GroupPut writes the given Group.
func (s *s3BucketStore) GroupPut(group *storagepb.Group) error {
	richGroup, err := group.ToRichGroup()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(richGroup, "", "\t")
	if err != nil {
		return err
	}

	return s.client.writeObject("groups/", group.Id+".json", data)
}

// GroupGet returns a machine Group by id.
func (s *s3BucketStore) GroupGet(id string) (*storagepb.Group, error) {
	data, err := s.client.readObject("groups/", id+".json")
	if err != nil {
		return nil, err
	}
	group, err := storagepb.ParseGroup(data)
	if err != nil {
		return nil, err
	}
	return group, err
}

// GroupDelete deletes a machine Group by id.
func (s *s3BucketStore) GroupDelete(id string) error {
	return s.client.deleteObject("groups/", id+".json")
}

// GroupList lists all machine Groups.
func (s *s3BucketStore) GroupList() ([]*storagepb.Group, error) {
	objects, err := s.client.listPrefix("groups/")
	if err != nil {
		return nil, err
	}
	groups := make([]*storagepb.Group, 0, len(objects))
	for _, obj := range objects {
		name := strings.TrimPrefix(*obj.Key, "groups/")
		name = strings.TrimSuffix(name, ".json")
		group, err := s.GroupGet(name)
		if err == nil {
			groups = append(groups, group)
		} else if s.logger != nil {
			s.logger.Infof("Group %q: %v", name, err)
		}
	}
	return groups, nil
}

// ProfilePut writes the given Profile.
func (s *s3BucketStore) ProfilePut(profile *storagepb.Profile) error {
	data, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return err
	}
	return s.client.writeObject("profiles/", profile.Id+".json", data)
}

// ProfileGet gets a profile by id.
func (s *s3BucketStore) ProfileGet(id string) (*storagepb.Profile, error) {
	data, err := s.client.readObject("profiles/", id+".json")
	if err != nil {
		return nil, err
	}
	profile := new(storagepb.Profile)
	err = json.Unmarshal(data, profile)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, err
}

// ProfileDelete deletes a profile by id.
func (s *s3BucketStore) ProfileDelete(id string) error {
	return s.client.deleteObject("profiles/", id+".json")
}

// ProfileList lists all profiles.
func (s *s3BucketStore) ProfileList() ([]*storagepb.Profile, error) {
	objects, err := s.client.listPrefix("profiles/")
	if err != nil {
		return nil, err
	}
	profiles := make([]*storagepb.Profile, 0, len(objects))
	for _, obj := range objects {
		name := strings.TrimPrefix(*obj.Key, "profiles/")
		name = strings.TrimSuffix(name, ".json")
		profile, err := s.ProfileGet(name)
		if err == nil {
			profiles = append(profiles, profile)
		} else if s.logger != nil {
			s.logger.Infof("Profile %q: %v", name, err)
		}
	}
	return profiles, nil
}

// IgnitionPut creates or updates an Ignition template.
func (s *s3BucketStore) IgnitionPut(name string, config []byte) error {
	return s.client.writeObject("ignition/", name, config)
}

// IgnitionGet gets an Ignition template by name.
func (s *s3BucketStore) IgnitionGet(name string) (string, error) {
	data, err := s.client.readObject("ignition/", name)
	return string(data), err
}

// IgnitionDelete deletes an Ignition template by name.
func (s *s3BucketStore) IgnitionDelete(name string) error {
	return s.client.deleteObject("ignition/", name)
}

// GenericPut creates or updates an Generic template.
func (s *s3BucketStore) GenericPut(name string, config []byte) error {
	return s.client.writeObject("generic/", name, config)
}

// GenericGet gets an Generic template by name.
func (s *s3BucketStore) GenericGet(name string) (string, error) {
	data, err := s.client.readObject("generic/", name)
	return string(data), err
}

// GenericDelete deletes an Generic template by name.
func (s *s3BucketStore) GenericDelete(name string) error {
	return s.client.deleteObject("generic/", name)
}

// CloudGet gets a Cloud-Config template by name.
func (s *s3BucketStore) CloudGet(name string) (string, error) {
	data, err := s.client.readObject("cloud/", name)
	return string(data), err
}
