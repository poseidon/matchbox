package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
	fake "github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestS3GroupGet(t *testing.T) {

	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}

	// Test error response
	_, err := store.GroupGet("blah")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientGet)
	}

	// Test Groups can be retrieved from bucket
	testClient.svc = fake.WorkingS3Client
	group, err := store.GroupGet("test-group")
	assert.Nil(t, err)
	assert.Equal(t, fake.Group, group)

	group, err = store.GroupGet("group-no-metadata")
	assert.Nil(t, err)
	assert.Equal(t, fake.GroupNoMetadata, group)
}

func TestS3GroupList(t *testing.T) {

	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}
	_, err := store.GroupList()
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientList)
	}

	testClient.svc = fake.WorkingS3Client
	groups, err := store.GroupList()
	assert.Nil(t, err)
	if assert.Equal(t, 2, len(groups)) {
		assert.Contains(t, groups, fake.Group)
		assert.Contains(t, groups, fake.GroupNoMetadata)
		assert.NotContains(t, groups, &storagepb.Group{})
	}
}

func TestS3ProfileGet(t *testing.T) {

	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}
	_, err := store.ProfileGet("blah")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientGet)
	}

	testClient.svc = fake.WorkingS3Client
	profile, err := store.ProfileGet(fake.Profile.Id)
	assert.Equal(t, fake.Profile, profile)
	assert.Nil(t, err)
}

func TestS3ProfileList(t *testing.T) {
	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}
	_, err := store.ProfileList()
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientList)
	}

	testClient.svc = fake.WorkingS3Client
	profiles, err := store.ProfileList()
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(profiles)) {
		assert.Equal(t, fake.Profile, profiles[0])
	}
}

func TestS3IgnitionGet(t *testing.T) {
	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}
	_, err := store.IgnitionGet("myignition.json")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientGet)
	}

	testClient.svc = fake.WorkingS3Client
	ign, err := store.IgnitionGet("myignition.json")
	assert.Equal(t, fake.IgnitionYAML, ign)
	assert.Nil(t, err)
}

func TestS3GenericGet(t *testing.T) {
	testClient := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}

	store := &s3BucketStore{
		client: testClient,
	}
	_, err := store.GenericGet("generic.json")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientGet)
	}

	testClient.svc = fake.WorkingS3Client
	gn, err := store.GenericGet("generic.json")
	assert.Equal(t, fake.Generic, gn)
	assert.Nil(t, err)
}
