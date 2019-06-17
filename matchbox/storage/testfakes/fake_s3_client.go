package testfakes

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	//"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

var (
	s3MockBucket = "s3MockBucket"

	// ErrMockClientGet to fake s3 client get errors
	ErrMockClientGet = errors.New("ErrClientGet")
	// ErrMockClientPut to fake s3 client put errors
	ErrMockClientPut = errors.New("ErrClientPut")
	// ErrMockClientDelete to fake s3 client delete errors
	ErrMockClientDelete = errors.New("ErrClientDelete")
	// ErrMockClientWaitDelete to fake s3 client wait on deletion errors
	ErrMockClientWaitDelete = errors.New("ErrClientWaitDelete")
	// ErrMockClientList to fake s3 client list errors
	ErrMockClientList = errors.New("ErrClientList")

	emptyGetObjectOutput    = &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}
	emptyPutObjectOutput    = &s3.PutObjectOutput{}
	emptyDeleteObjectOutput = &s3.DeleteObjectOutput{}
	emptyListObjectsOutput  = &s3.ListObjectsV2Output{}
)

type mockS3Client struct {
	s3iface.S3API
	getObjectMock                func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	putObjectMock                func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	deleteObjectMock             func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	waitUntilObjectNotExistsMock func(*s3.HeadObjectInput) error
	listObjectsV2Mock            func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
}

func (m mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return m.putObjectMock(input)
}

func (m mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.getObjectMock(input)
}

func (m mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return m.deleteObjectMock(input)
}

func (m mockS3Client) WaitUntilObjectNotExists(input *s3.HeadObjectInput) error {
	return m.waitUntilObjectNotExistsMock(input)
}

func (m mockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return m.listObjectsV2Mock(input)
}

// ErrorS3Client s3 clients that always return an error
var ErrorS3Client = mockS3Client{
	getObjectMock: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		return emptyGetObjectOutput, ErrMockClientGet
	},
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return emptyPutObjectOutput, ErrMockClientPut
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return emptyDeleteObjectOutput, ErrMockClientDelete
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return ErrMockClientWaitDelete
	},
	listObjectsV2Mock: func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return emptyListObjectsOutput, ErrMockClientList
	},
}

// NoErrorEmptyS3Client mocks successful responses with empty content
var NoErrorEmptyS3Client = mockS3Client{
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return emptyPutObjectOutput, nil
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return emptyDeleteObjectOutput, nil
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return nil
	},
	listObjectsV2Mock: func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return emptyListObjectsOutput, nil
	},
}

// WorkingS3Client mocks a working S3client
var WorkingS3Client = mockS3Client{
	getObjectMock: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if strings.HasPrefix(*input.Key, "groups") {
			return returnGroup(*input.Key)
		}
		if strings.HasPrefix(*input.Key, "profiles") {
			return returnProfile(*input.Key)
		}
		if strings.HasPrefix(*input.Key, "ignition") {
			return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewBufferString(IgnitionYAML))}, nil
		}
		if strings.HasPrefix(*input.Key, "generic") {
			return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewBufferString(Generic))}, nil
		}
		return emptyGetObjectOutput, ErrMockClientGet
	},
	putObjectMock: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return emptyPutObjectOutput, ErrMockClientPut
	},
	deleteObjectMock: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
		return emptyDeleteObjectOutput, ErrMockClientDelete
	},
	waitUntilObjectNotExistsMock: func(*s3.HeadObjectInput) error {
		return ErrMockClientWaitDelete
	},
	listObjectsV2Mock: func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		if *input.Prefix == "groups/" {
			list := []*s3.Object{
				&s3.Object{Key: aws.String("groups/" + Group.Id + ".json")},
				&s3.Object{Key: aws.String("groups/" + GroupNoMetadata.Id + ".json")},
			}
			return &s3.ListObjectsV2Output{Contents: list}, nil
		}
		if *input.Prefix == "profiles/" {
			list := []*s3.Object{
				&s3.Object{Key: aws.String("profiles/" + Profile.Id + ".json")},
			}
			return &s3.ListObjectsV2Output{Contents: list}, nil
		}
		return emptyListObjectsOutput, ErrMockClientList
	},
}

func returnGroup(gID string) (*s3.GetObjectOutput, error) {

	if gID == *aws.String("groups/" + Group.Id + ".json") {
		//json marshal will give base64 encoding for matadata
		raw := []byte(`
{
	"id": "test-group",
	"Name": "test group",
	"profile": "g1h2i3j4",
	"selector": {
		"uuid": "a1b2c3d4"
	},
	"metadata": {
		"pod_network": "10.2.0.0/16",
		"service_name":"etcd2"
	}
}
`)
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(raw))}, nil
	}
	if gID == *aws.String("groups/" + GroupNoMetadata.Id + ".json") {
		data, err := json.MarshalIndent(GroupNoMetadata, "", "\t")
		if err != nil {
			return emptyGetObjectOutput, err
		}
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(data))}, nil
	}
	return emptyGetObjectOutput, ErrMockClientGet
}

func returnProfile(pID string) (*s3.GetObjectOutput, error) {
	if pID == *aws.String("profiles/" + Profile.Id + ".json") {
		data, err := json.MarshalIndent(Profile, "", "\t")
		if err != nil {
			return emptyGetObjectOutput, err
		}
		return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(data))}, nil
	}
	return emptyGetObjectOutput, ErrMockClientGet
}
