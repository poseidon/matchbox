package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	fake "github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestReadObject(t *testing.T) {

	// Test no error when object exists
	s3Client := &S3Client{
		svc:    fake.WorkingS3Client,
		bucket: "test-bucket",
	}

	_, err := s3Client.readObject("groups/", fake.Group.Id+".json")
	assert.Nil(t, err)

	// Test that we get the client error when the request fails
	s3Client.svc = fake.ErrorS3Client
	_, err = s3Client.readObject("groups/", "blah.json")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientGet)
	}
}

func TestWriteObject(t *testing.T) {
	// Test no error is returned when s3 client request succeeds
	s3Client := &S3Client{
		svc:    fake.NoErrorEmptyS3Client,
		bucket: "test-bucket",
	}

	data := []byte{}

	err := s3Client.writeObject("groups/", "test.json", data)
	assert.Nil(t, err)

	// Test that we get the client error when the request fails
	s3Client.svc = fake.ErrorS3Client
	err = s3Client.writeObject("groups/", "test.json", data)
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientPut)
	}
}

func TestDeleteObject(t *testing.T) {
	s3Client := &S3Client{
		svc:    fake.NoErrorEmptyS3Client,
		bucket: "test-bucket",
	}

	// Test no error is returned when s3 client request succeeds
	err := s3Client.deleteObject("groups/", "test.json")
	assert.Nil(t, err)

	// Test that we get the client error when the request fails
	s3Client.svc = fake.ErrorS3Client
	err = s3Client.deleteObject("groups/", "blah.json")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientDelete)
	}
}

func TestListPrefix(t *testing.T) {
	// Test error piped when client returns error
	s3Client := &S3Client{
		svc:    fake.ErrorS3Client,
		bucket: "test-bucket",
	}
	_, err := s3Client.listPrefix("groups/")
	if assert.Error(t, err) {
		assert.Error(t, err, fake.ErrMockClientList)
	}

	// Test empty list
	s3Client.svc = fake.NoErrorEmptyS3Client
	list, err := s3Client.listPrefix("groups/")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	// Test group list with 2 groups
	s3Client.svc = fake.WorkingS3Client
	list, err = s3Client.listPrefix("groups/")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
}
